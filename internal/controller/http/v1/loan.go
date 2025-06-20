package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	httpHelper "github.com/ferdikurniawan/loan-service/internal/controller/http"
	"github.com/ferdikurniawan/loan-service/internal/entity"
	"github.com/ferdikurniawan/loan-service/internal/services"
)

type loanRoutes struct {
	loanService services.LoanService
}

func newLoanRoutes(handler *gin.RouterGroup, svc services.LoanService) {
	r := &loanRoutes{svc}

	handler.POST("/loans", r.submitLoan)                      //borrower submits a new Loan
	handler.PATCH("/loans/:loan_id/status", r.updateLoan)     //update Loan status
	handler.POST("/loans/:loan_id/investments", r.investLoan) //investor chip in
	handler.POST("/loans/:loan_id/disburse", r.disburseLoan)  //disbursement
	handler.GET("/loans/:loan_id", r.getLoan)                 //get loan detail
}

func (r *loanRoutes) submitLoan(c *gin.Context) {

	borrowerID := c.GetString("borrowerID")
	var req entity.LoanSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value"},
			nil,
			http.StatusBadRequest,
		)
		return
	}
	req.BorrowerID = borrowerID //attach borrowerID to the request object since it does not exist in request payload, but obtained through context

	loan, err := r.loanService.CreateLoan(c, req)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: err.Error()},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	httpHelper.Response(c,
		true,
		nil,
		loan,
		http.StatusOK,
	)
}

func (r *loanRoutes) updateLoan(c *gin.Context) {

	staffID := c.GetString("staffID")
	loanID := c.Param("loan_id")

	if loanID == "" {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value (loan ID)"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	var req entity.LoanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	req.LoanID = loanID
	req.StaffID = staffID

	err := r.loanService.UpdateLoan(c, req)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: err.Error()},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	httpHelper.Response(c,
		true,
		nil,
		nil,
		http.StatusOK,
	)
}

func (r *loanRoutes) investLoan(c *gin.Context) {

	investorID := c.GetString("investorID")
	loanID := c.Param("loan_id")

	if loanID == "" {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value (loan ID)"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	var req entity.LoanInvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	req.LoanID = loanID
	req.InvestorID = uuid.MustParse(investorID)

	err := r.loanService.InvestLoan(c, req)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: err.Error()},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	httpHelper.Response(c,
		true,
		nil,
		nil,
		http.StatusOK,
	)

}

func (r *loanRoutes) disburseLoan(c *gin.Context) {

	staffID := c.GetString("staffID")
	loanID := c.Param("loan_id")

	if loanID == "" {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value (loan ID)"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	var req entity.LoanDisburseRequest
	if err := c.ShouldBind(&req); err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	//processing file: loan agreement signed by borrower
	err := c.Request.ParseMultipartForm(10 << 20) //10MB upper limit for memory usage
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "File is too large"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	disbursementDate := c.PostForm("disbursement_date")
	disburseAt, err := time.Parse("2006-01-02", disbursementDate)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Disbursement date cannot be empty / invalid"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	fileHeader, err := c.FormFile("agreement_file")
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing agreement file"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Failed to open agreement file"},
			nil,
			http.StatusBadRequest,
		)
		return
	}
	defer file.Close()

	// Save file to disk first (based on requirement, this should be uploaded to S3 but this is for demo purpose)
	savePath := fmt.Sprintf("./uploads/agreements/%s_%s", loanID, fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: "Failed to store agreement file"},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	req.LoanID = loanID
	req.DisburseAt = disburseAt
	req.StaffID = staffID
	req.AgreementLetterLink = savePath

	err = r.loanService.DisburseLoan(c, req)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: err.Error()},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	httpHelper.Response(c,
		true,
		nil,
		nil,
		http.StatusOK,
	)
}

func (r *loanRoutes) getLoan(c *gin.Context) {

	loanID := c.Param("loan_id")

	if loanID == "" {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "Missing / invalid required value (loan ID)"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	//loanID must be UUID
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 400, Type: "bad_request", Message: "invalid required value (loan ID)"},
			nil,
			http.StatusBadRequest,
		)
		return
	}

	loan, err := r.loanService.GetLoanByID(c, loanUUID)
	if err != nil {
		httpHelper.Response(c,
			false,
			&entity.ErrorResponse{Code: 500, Type: "server_error", Message: "internal server error"},
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	httpHelper.Response(c,
		true,
		nil,
		loan,
		http.StatusOK,
	)
}
