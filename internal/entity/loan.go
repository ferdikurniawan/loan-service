package entity

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	ID              uuid.UUID `json:"loan_id"`
	BorrowerID      uuid.UUID `json:"borrower_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestRate    float32   `json:"interest_rate"`
	AgreementLetter string    `json:"agreement_letter"`
	PublicID        string    `json:"loan_pub_id"`
	Status          string    `json:"loan_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DisburseAt      time.Time `json:"disburse_at"`
}

type LoanInvestment struct {
	ID         uuid.UUID `json:"loan_investment_id"`
	LoanID     uuid.UUID `json:"loan_id"`
	InvestorID uuid.UUID `json:"investor_id"`
	Amount     float64   `json:"amount"`
	InvestedAt time.Time `json:"invested_at"`
}

type LoanSubmitRequest struct {
	BorrowerID      string  `json:"-"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestRate    float32 `json:"interest_rate"`
	Reason          string  `json:"reason"`
}

type LoanUpdateRequest struct {
	LoanID  string `json:"-"`
	Status  string `json:"status"`
	StaffID int64  `json:"-"`
}

type LoanInvestRequest struct {
	LoanID     string    `json:"-"`
	Amount     float64   `json:"amount"`
	InvestorID uuid.UUID `json:"-"`
}

type LoanDisburseRequest struct {
	LoanID              string
	BorrowerID          int64
	StaffID             int64
	LoanAgreementDocs   *multipart.FileHeader `form:"agreement_file" binding:"required"`
	DisbursementDate    string                `form:"disbursement_date" binding:"required"`
	AgreementLetterLink string
	DisburseAt          time.Time
}
