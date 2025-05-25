package http

import (
	"github.com/ferdikurniawan/loan-service/internal/entity"
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, success bool, err entity.ErrorResponse, data interface{}, httpStatus int) {
	c.JSON(httpStatus, gin.H{
		"success": success,
		"error":   err,
		"data":    data,
	})
}
