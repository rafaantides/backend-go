package middlewares

import (
	"backend-go/internal/api/errs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Valida parâmetros de consulta da URL.
		for key, values := range c.Request.URL.Query() {
			if strings.Contains(key, "id") {
				for _, value := range values {
					_, err := uuid.Parse(value)
					if err != nil {
						// TODO: rever aqui
						c.Error(errs.NewAPIError(http.StatusBadRequest, errs.AAAInvalidUUID(key, err)))
						formatError(c)
						c.Abort()
						return
					}
				}
			}
		}

		// Valida parâmetros do corpo da requisição (JSON).
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(errs.NewAPIError(http.StatusBadRequest, err))
			formatError(c)
			c.Abort()
			return

		}

		if err := validateUUIDsRecursive(body); err != nil {
			c.Error(errs.NewAPIError(http.StatusBadRequest, err))
			formatError(c)
			c.Abort()
			return
		}

		c.Next()
	}

}

func validateUUIDsRecursive(data map[string]interface{}) error {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if strings.Contains(key, "id") {
				_, err := uuid.Parse(v)
				if err != nil {
					return errs.AAAInvalidUUID(key, err)
				}
			}
		case []any:
			for _, item := range v {
				if strItem, ok := item.(string); ok {
					_, err := uuid.Parse(strItem)
					if err != nil {
						return errs.AAAInvalidUUID(key, err)
					}
				}
			}
		case map[string]any:
			if err := validateUUIDsRecursive(v); err != nil {
				return err
			}
		}
	}
	return nil
}
