package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/ferdikurniawan/loan-service/config"
	"github.com/ferdikurniawan/loan-service/internal/services"
	"github.com/ferdikurniawan/loan-service/internal/utils"
)

type Services struct {
	Cfg *config.Config

	LoanService services.LoanService
}

func (s Services) Initialized() error {
	return utils.ValidateStruct(s)
}

func NewRouter(handler *gin.Engine, s Services) {

	// panic if any of the services field is not initialized
	if err := s.Initialized(); err != nil {
		panic(err)
	}

	// Routers
	h := handler.Group("v1")
	{
		newLoanRoutes(h, s.LoanService) //TODO use custom middleware for AUTH
	}
}
