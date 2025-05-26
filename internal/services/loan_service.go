package services

import (
	"context"
	"errors"
	"log"

	"github.com/ferdikurniawan/loan-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=loan_service.go -package=mock -destination=mock/loan_service_mock.go
type (
	LoanService interface {
		CreateLoan(ctx context.Context, loanRequest entity.LoanSubmitRequest) (*entity.Loan, error)
		UpdateLoan(ctx context.Context, loanStatusRequest entity.LoanUpdateRequest) error
		InvestLoan(ctx context.Context, loanInvestRequest entity.LoanInvestRequest) error
		DisburseLoan(ctx context.Context, loanDisburseRequest entity.LoanDisburseRequest) error
		GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error)
	}

	loanService struct {
		repo LoanRepo
	}

	LoanRepo interface {
		InsertLoan(ctx context.Context, loan *entity.Loan) (*entity.Loan, error)
		UpdateLoanStatus(ctx context.Context, loan *entity.Loan, staffID uuid.UUID) error
		AddLoanInvestments(ctx context.Context, investment entity.LoanInvestment) error
		DisburseLoan(ctx context.Context, loan *entity.Loan, staffID uuid.UUID) error
		GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error)
	}
)

func NewLoanService(repo LoanRepo) *loanService {
	return &loanService{
		repo: repo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, loanRequest entity.LoanSubmitRequest) (*entity.Loan, error) {

	loan := entity.Loan{
		ID:              uuid.New(),
		BorrowerID:      uuid.MustParse(loanRequest.BorrowerID),
		PrincipalAmount: loanRequest.PrincipalAmount,
		InterestRate:    loanRequest.InterestRate,
	}

	res, err := s.repo.InsertLoan(ctx, &loan)
	if err != nil {
		log.Printf("[CreateLoan] error creating loan: %s", err.Error())
	}

	return res, err
}

func (s *loanService) UpdateLoan(ctx context.Context, loanStatusRequest entity.LoanUpdateRequest) error {

	loan := entity.Loan{
		ID:     uuid.MustParse(loanStatusRequest.LoanID),
		Status: loanStatusRequest.Status,
	}

	err := s.repo.UpdateLoanStatus(ctx, &loan, uuid.MustParse(loanStatusRequest.StaffID))
	if err != nil {
		log.Printf("[UpdateLoan] error update loan: %s", err.Error())
	}
	return err
}

func (s *loanService) InvestLoan(ctx context.Context, loanInvestRequest entity.LoanInvestRequest) error {

	investment := entity.LoanInvestment{
		LoanID:     uuid.MustParse(loanInvestRequest.LoanID),
		Amount:     loanInvestRequest.Amount,
		InvestorID: loanInvestRequest.InvestorID,
	}

	err := s.repo.AddLoanInvestments(ctx, investment)
	if err != nil {
		log.Printf("[InvestLoan] error invest loan: %s", err.Error())
	}

	return err
}

func (s *loanService) GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	loan.Returns = float64(loan.InterestRate/100) * float64(loan.PrincipalAmount)
	return loan, err

}

func (s *loanService) DisburseLoan(ctx context.Context, loanDisburseRequest entity.LoanDisburseRequest) error {

	loan := entity.Loan{
		ID:              uuid.MustParse(loanDisburseRequest.LoanID),
		Status:          "disbursed",
		AgreementLetter: loanDisburseRequest.AgreementLetterLink,
		DisburseAt:      loanDisburseRequest.DisburseAt,
	}

	currentLoan, err := s.GetLoanByID(ctx, uuid.MustParse(loanDisburseRequest.LoanID))
	if err != nil {
		log.Printf("[DisburseLoan] error getting loan detail: %s", err.Error())
		return err
	}

	if currentLoan.Status != "invested" {
		return errors.New("loan principal amount is not met yet")
	}

	err = s.repo.DisburseLoan(ctx, &loan, uuid.MustParse(loanDisburseRequest.StaffID))
	if err != nil {
		log.Printf("[DisburseLoan] error disburse loan: %s", err.Error())
	}
	return err
}
