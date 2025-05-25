package services

import (
	"context"
	"errors"

	"github.com/ferdikurniawan/loan-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=loan_service.go -package=mock -destination=mock/loan_service_mock.go
type (
	LoanService interface {
		CreateLoan(ctx context.Context, loanRequest entity.LoanSubmitRequest) error
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
		UpdateLoanStatus(ctx context.Context, loan *entity.Loan, staffID int64) error
		AddLoanInvestments(ctx context.Context, investment entity.LoanInvestment) error
		DisburseLoan(ctx context.Context, loan *entity.Loan, staffID int64) error
		GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error)
	}
)

func NewLoanService(repo LoanRepo) *loanService {
	return &loanService{
		repo: repo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, loanRequest entity.LoanSubmitRequest) error {

	loan := entity.Loan{
		BorrowerID:      loanRequest.BorrowerID,
		PrincipalAmount: loanRequest.PrincipalAmount,
		InterestRate:    loanRequest.InterestRate,
		PublicID:        "public-id", //TODO use ULID

	}

	_, err := s.repo.InsertLoan(ctx, &loan)

	return err
}

func (s *loanService) UpdateLoan(ctx context.Context, loanStatusRequest entity.LoanUpdateRequest) error {

	loan := entity.Loan{
		ID:     uuid.MustParse(loanStatusRequest.LoanID),
		Status: loanStatusRequest.Status,
	}

	err := s.repo.UpdateLoanStatus(ctx, &loan, loanStatusRequest.StaffID)
	return err
}

func (s *loanService) InvestLoan(ctx context.Context, loanInvestRequest entity.LoanInvestRequest) error {

	investment := entity.LoanInvestment{
		LoanID:     uuid.MustParse(loanInvestRequest.LoanID),
		Amount:     loanInvestRequest.Amount,
		InvestorID: loanInvestRequest.InvestorID,
	}

	err := s.repo.AddLoanInvestments(ctx, investment)
	//TODO send email stub if status is invested
	return err
}

func (s *loanService) GetLoanByID(ctx context.Context, loanID uuid.UUID) (*entity.Loan, error) {
	loan, err := s.repo.GetLoanByID(ctx, loanID)
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
		return err
	}

	if currentLoan.Status != "invested" {
		return errors.New("loan principal amount is not met yet")
	}

	err = s.repo.DisburseLoan(ctx, &loan, loanDisburseRequest.StaffID)
	return err
}
