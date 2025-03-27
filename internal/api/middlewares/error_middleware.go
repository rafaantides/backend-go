package middlewares

import (
	"backend-go/internal/api/errs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		formatError(c)

	}
}

func formatError(c *gin.Context) {
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			log.Printf("[ERROR] %s", err.Err)
		}

		err := c.Errors.Last().Err

		if apiErr, ok := err.(*errs.APIError); ok {
			statusCode := apiErr.StatusCode
			message, exists := errs.ErrorMessages[statusCode]
			if !exists {
				message = errs.InternalServerError
			}

			c.JSON(statusCode, errs.ErrorResponse{
				Message: message,
				Detail:  err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
				Message: "Ocorreu um erro inesperado",
				Detail:  err.Error(),
			})
		}
	}
}
