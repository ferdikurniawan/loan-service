package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ferdikurniawan/loan-service/internal/pkg/postgres"
	"github.com/google/uuid"

	"github.com/ferdikurniawan/loan-service/internal/entity"
)

type (
	loanRepo struct {
		*postgres.Postgres
	}
)

func NewLoanRepo(pg *postgres.Postgres) *loanRepo {
	return &loanRepo{pg}
}

func (r *loanRepo) InsertLoan(ctx context.Context, loan *entity.Loan) (*entity.Loan, error) {
	var result entity.Loan

	query := `INSERT INTO loan (loan_id, borrower_id, principal_amount, interest_rate, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING loan_id, created_at, updated_at`
	err := r.DB.QueryRowContext(ctx, query, loan.ID, loan.BorrowerID, loan.PrincipalAmount, loan.InterestRate, "proposed", "now()", "now()").Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}

	result.BorrowerID = loan.BorrowerID
	result.PrincipalAmount = loan.PrincipalAmount
	result.InterestRate = loan.InterestRate
	result.Status = "proposed"

	return &result, nil
}

func (r *loanRepo) UpdateLoanStatus(ctx context.Context, loan *entity.Loan, staffID uuid.UUID) error {

	//wrap queries within one transaction since there are multiple dependent ops:
	//update status and loan log table record insertion
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var updatedAt time.Time
	var currentStatus string
	query := `SELECT updated_at, status FROM loan WHERE loan_id = $1` //get loan detail, esp the updated_at to achieve optimistic locking
	err = tx.QueryRowContext(ctx, query, loan.ID).Scan(&updatedAt, &currentStatus)
	if err == sql.ErrNoRows {
		return fmt.Errorf("loan not found")
	} else if err != nil {
		return err
	}

	updateTime := time.Now()
	queryUpdate := `UPDATE loan SET status = $4, updated_at = $3 WHERE loan_id = $1 AND updated_at = $2`

	res, err := tx.ExecContext(ctx, queryUpdate, loan.ID, updatedAt, updateTime, loan.Status)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("loan has been updated by another staff")
	}

	loanPrev := entity.Loan{
		ID:        loan.ID,
		Status:    currentStatus,
		UpdatedAt: updatedAt,
	}
	loanAfter := loanPrev
	loanAfter.Status = loan.Status
	loanAfter.UpdatedAt = updateTime

	queryLoanStatusHistory := `INSERT INTO loan_status_history (loan_id, before, after, updated_by, updated_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(queryLoanStatusHistory, loan.ID, loanPrev, loanAfter, staffID, "now()")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *loanRepo) AddLoanInvestments(ctx context.Context, investment entity.LoanInvestment) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//1. Check amount & status of the loan
	var amount int64
	var status string
	query := `SELECT principal_amount, status FROM loan where loan_id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, investment.LoanID).Scan(&amount, &status)
	if err != nil {
		return err
	}
	if status != "approved" {
		return errors.New("loan is not approved yet / has reach principal amount")
	}

	//2. Get total investment
	var totalInvested int64
	query = `SELECT COALESCE(SUM(amount), 0) FROM loan_investment where loan_id = $1`
	err = tx.QueryRowContext(ctx, query, investment.LoanID).Scan(&totalInvested)
	if err != nil {
		return err
	}

	remaining := amount - totalInvested
	if investment.Amount > remaining {
		return errors.New("pledged fund exceeds the remaining loan value")
	}

	//3. Insert the investment
	query = `INSERT INTO loan_investment (loan_investment_id, loan_id, investor_id, amount, invested_at)
	VALUES ($1, $2, $3, $4, 'now()')`
	_, err = tx.ExecContext(ctx, query, uuid.New(), investment.LoanID, investment.InvestorID, investment.Amount)
	if err != nil {
		return err
	}

	//4. Update Loan status if invested fund reached principal loan amount
	if investment.Amount == remaining {
		updatedTime := time.Now()
		query = `UPDATE loan SET status = 'invested', updated_at = $2 WHERE loan_id = $1`
		_, err = tx.ExecContext(ctx, query, investment.LoanID, updatedTime)
		if err != nil {
			return err
		}

		//4.1 Add Loan History Log Record
		loanPrev := entity.Loan{
			ID:     investment.LoanID,
			Status: status,
		}
		loanAfter := loanPrev
		loanAfter.Status = "invested"
		loanAfter.UpdatedAt = updatedTime

		queryLoanStatusHistory := `INSERT INTO loan_status_history (loan_id, before, after, updated_at)
		VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(queryLoanStatusHistory, investment.LoanID, loanPrev, loanAfter, "now()")
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}

func (r *loanRepo) GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error) {

	var (
		loan            entity.Loan
		agreementLetter sql.NullString
		updatedAt       sql.NullTime
		disburseAt      sql.NullTime
	)

	query := `SELECT loan_id, borrower_id, principal_amount, interest_rate, agreement_letter, status, created_at, updated_at, disburse_at
	FROM loan WHERE loan_id = $1`
	err := r.DB.QueryRowContext(ctx, query, loanID).Scan(&loan.ID, &loan.BorrowerID,
		&loan.PrincipalAmount, &loan.InterestRate, &agreementLetter, &loan.Status, &loan.CreatedAt,
		&updatedAt, &disburseAt)

	if err != nil {
		return nil, err
	}

	loan.AgreementLetter = agreementLetter.String
	loan.UpdatedAt = updatedAt.Time
	loan.DisburseAt = disburseAt.Time

	return &loan, err

}

func (r *loanRepo) DisburseLoan(ctx context.Context, loan *entity.Loan, staffID uuid.UUID) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE loan SET status = $1, agreement_letter = $2, disburse_at = $3, updated_at = $5 WHERE loan_id = $4`
	_, err = tx.ExecContext(ctx, query, loan.Status, loan.AgreementLetter, loan.DisburseAt, loan.ID, "now()")
	if err != nil {
		return err
	}

	loanPrev := entity.Loan{
		ID:              loan.ID,
		Status:          "invested",
		AgreementLetter: "",
	}
	loanAfter := loanPrev
	loanAfter.Status = loan.Status
	loanAfter.AgreementLetter = loan.AgreementLetter

	queryLoanStatusHistory := `INSERT INTO loan_status_history (loan_id, before, after, updated_by, updated_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(queryLoanStatusHistory, loan.ID, loanPrev, loanAfter, staffID, "now()")
	if err != nil {
		return err
	}

	return tx.Commit()
}
