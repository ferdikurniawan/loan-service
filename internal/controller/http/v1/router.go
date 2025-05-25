package v1

import (
	"net/http"

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
	h.Use(DummyAuthMiddleware())
	{
		newLoanRoutes(h, s.LoanService) //TODO use custom middleware for AUTH
	}
}

// Dummy auth to get the app working, we can differentiate actor based on Header
func DummyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		staffID := c.GetHeader("X-Staff-ID")
		borrowerID := c.GetHeader("X-Borrower-ID")

		if staffID == "" && borrowerID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing X-Staff-ID or X-Borrower-ID",
			})
			return
		}

		if staffID != "" {
			c.Set("staffID", staffID)
		}
		if borrowerID != "" {
			c.Set("borrowerID", borrowerID)
		}

		c.Next()
	}
}
