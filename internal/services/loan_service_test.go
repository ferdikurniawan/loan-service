package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ferdikurniawan/loan-service/internal/entity"
	mock "github.com/ferdikurniawan/loan-service/internal/services/mock"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupLoanService(t *testing.T) (*loanService, *mock.MockLoanRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockLoanRepo(ctrl)

	svc := NewLoanService(repo)

	return svc, repo
}

func Test_CreateLoan(t *testing.T) {
	t.Parallel()

	svc, repo := setupLoanService(t)
	ctx := context.Background()

	t.Run("create loan failed, error insert loan data to DB", func(t *testing.T) {

		loanReq := entity.LoanSubmitRequest{
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Reason:          "business reason",
			BorrowerID:      "d149aaa5-e7e8-4820-93a0-e278dcde447a",
		}

		repo.EXPECT().InsertLoan(ctx, gomock.Any()).Return(nil, errors.New("error db"))
		_, err := svc.CreateLoan(ctx, loanReq)
		assert.Equal(t, err.Error(), "error db")
	})
	t.Run("create loan success", func(t *testing.T) {

		loanReq := entity.LoanSubmitRequest{
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Reason:          "business reason",
			BorrowerID:      "d149aaa5-e7e8-4820-93a0-e278dcde447a",
		}

		loanData := entity.Loan{
			ID:              uuid.MustParse("a98ba4bd-1e09-4134-b244-d89f9a86a44c"),
			BorrowerID:      uuid.MustParse("d149aaa5-e7e8-4820-93a0-e278dcde447a"),
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Status:          "proposed",
		}

		repo.EXPECT().InsertLoan(ctx, gomock.Any()).Return(&loanData, nil)
		loan, err := svc.CreateLoan(ctx, loanReq)
		assert.Nil(t, err)
		assert.Equal(t, loan.ID, loanData.ID)
	})
}

func Test_UpdateLoan(t *testing.T) {
	t.Parallel()

	svc, repo := setupLoanService(t)
	ctx := context.Background()

	t.Run("upload loan status failed, error db", func(t *testing.T) {
		loan := entity.Loan{
			ID:     uuid.MustParse("3e6a779e-d857-4ad3-af95-693d16e6f6d1"),
			Status: "approved",
		}

		loanUpdateReq := entity.LoanUpdateRequest{
			LoanID:  "3e6a779e-d857-4ad3-af95-693d16e6f6d1",
			Status:  "approved",
			StaffID: "1e938a3c-3752-49a6-a2a6-43be38c6aa82",
		}

		repo.EXPECT().UpdateLoanStatus(ctx, &loan, uuid.MustParse("1e938a3c-3752-49a6-a2a6-43be38c6aa82")).Return(errors.New("db error"))

		err := svc.UpdateLoan(ctx, loanUpdateReq)
		assert.Equal(t, err.Error(), "db error")
	})

	t.Run("upload loan status success", func(t *testing.T) {
		loan := entity.Loan{
			ID:     uuid.MustParse("3e6a779e-d857-4ad3-af95-693d16e6f6d1"),
			Status: "approved",
		}

		loanUpdateReq := entity.LoanUpdateRequest{
			LoanID:  "3e6a779e-d857-4ad3-af95-693d16e6f6d1",
			Status:  "approved",
			StaffID: "1e938a3c-3752-49a6-a2a6-43be38c6aa82",
		}

		repo.EXPECT().UpdateLoanStatus(ctx, &loan, uuid.MustParse("1e938a3c-3752-49a6-a2a6-43be38c6aa82")).Return(nil)

		err := svc.UpdateLoan(ctx, loanUpdateReq)
		assert.Nil(t, err)
	})
}

func Test_InvestLoan(t *testing.T) {
	t.Parallel()

	svc, repo := setupLoanService(t)
	ctx := context.Background()

	t.Run("invest loan failed, error when adding records to db", func(t *testing.T) {

		loanInvestReq := entity.LoanInvestRequest{
			InvestorID: uuid.MustParse("e217fd14-0de2-4a11-8989-d8d51e2b9886"),
			LoanID:     "36b84065-1de5-47df-b1a6-311ff28dfe5b",
			Amount:     500000,
		}

		investment := entity.LoanInvestment{
			InvestorID: loanInvestReq.InvestorID,
			Amount:     loanInvestReq.Amount,
			LoanID:     uuid.MustParse(loanInvestReq.LoanID),
		}

		repo.EXPECT().AddLoanInvestments(ctx, investment).Return(errors.New("error adding db records"))

		err := svc.InvestLoan(ctx, loanInvestReq)
		assert.Equal(t, err.Error(), "error adding db records")
	})

	t.Run("invest loan success", func(t *testing.T) {

		loanInvestReq := entity.LoanInvestRequest{
			InvestorID: uuid.MustParse("e217fd14-0de2-4a11-8989-d8d51e2b9886"),
			LoanID:     "36b84065-1de5-47df-b1a6-311ff28dfe5b",
			Amount:     500000,
		}

		investment := entity.LoanInvestment{
			InvestorID: loanInvestReq.InvestorID,
			Amount:     loanInvestReq.Amount,
			LoanID:     uuid.MustParse(loanInvestReq.LoanID),
		}

		repo.EXPECT().AddLoanInvestments(ctx, investment).Return(nil)

		err := svc.InvestLoan(ctx, loanInvestReq)
		assert.Nil(t, err)
	})
}

func Test_DisburseLoan(t *testing.T) {
	t.Parallel()

	svc, repo := setupLoanService(t)
	ctx := context.Background()

	t.Run("disburse loan failed, error getting loan detail", func(t *testing.T) {
		loanDisburseReq := entity.LoanDisburseRequest{
			AgreementLetterLink: "./uploads/agreement.pdf",
			DisbursementDate:    "2025-05-25",
			LoanID:              "2badced4-3fa0-4a7e-8dcf-7c8031f0e704",
		}

		repo.EXPECT().GetLoanByID(ctx, uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704")).Return(nil, errors.New("error db"))

		err := svc.DisburseLoan(ctx, loanDisburseReq)
		assert.Equal(t, err.Error(), "error db")
	})

	t.Run("disburse loan failed, loan status is not invested yet", func(t *testing.T) {
		loanDisburseReq := entity.LoanDisburseRequest{
			AgreementLetterLink: "./uploads/agreement.pdf",
			DisbursementDate:    "2025-05-25",
			LoanID:              "2badced4-3fa0-4a7e-8dcf-7c8031f0e704",
		}

		loan := entity.Loan{
			ID:              uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704"),
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Status:          "approved",
		}

		repo.EXPECT().GetLoanByID(ctx, uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704")).Return(&loan, nil)

		err := svc.DisburseLoan(ctx, loanDisburseReq)
		assert.Equal(t, err.Error(), "loan principal amount is not met yet")
	})

	t.Run("disburse loan failed, error on db during disbursement", func(t *testing.T) {
		loanDisburseReq := entity.LoanDisburseRequest{
			AgreementLetterLink: "./uploads/agreement.pdf",
			DisbursementDate:    "2025-05-25",
			LoanID:              "2badced4-3fa0-4a7e-8dcf-7c8031f0e704",
			StaffID:             "75ed6802-8f18-4c5e-95b6-e8bd35e8d940",
		}

		loan := entity.Loan{
			ID:              uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704"),
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Status:          "invested",
		}

		repo.EXPECT().GetLoanByID(ctx, uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704")).Return(&loan, nil)
		repo.EXPECT().DisburseLoan(ctx, gomock.Any(), uuid.MustParse("75ed6802-8f18-4c5e-95b6-e8bd35e8d940")).Return(errors.New("db query error when disbursement"))

		err := svc.DisburseLoan(ctx, loanDisburseReq)
		assert.Equal(t, err.Error(), "db query error when disbursement")
	})

	t.Run("disburse loan success", func(t *testing.T) {
		loanDisburseReq := entity.LoanDisburseRequest{
			AgreementLetterLink: "./uploads/agreement.pdf",
			DisbursementDate:    "2025-05-25",
			LoanID:              "2badced4-3fa0-4a7e-8dcf-7c8031f0e704",
			StaffID:             "75ed6802-8f18-4c5e-95b6-e8bd35e8d940",
		}

		loan := entity.Loan{
			ID:              uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704"),
			PrincipalAmount: 1000000,
			InterestRate:    10.0,
			Status:          "invested",
		}

		repo.EXPECT().GetLoanByID(ctx, uuid.MustParse("2badced4-3fa0-4a7e-8dcf-7c8031f0e704")).Return(&loan, nil)
		repo.EXPECT().DisburseLoan(ctx, gomock.Any(), uuid.MustParse("75ed6802-8f18-4c5e-95b6-e8bd35e8d940")).Return(nil)

		err := svc.DisburseLoan(ctx, loanDisburseReq)
		assert.Nil(t, err)
	})
}
