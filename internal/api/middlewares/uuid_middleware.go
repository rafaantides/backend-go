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
						c.Error(errs.NewAPIError(http.StatusBadRequest, errs.InvalidParam(key, err)))
						handleError(c)
						c.Abort()
						return
					}
				}
			}
		}

		if err := validateBodyUUIDs(c); err != nil {
			c.Error(errs.NewAPIError(http.StatusBadRequest, err))
			handleError(c)
			c.Abort()
			return
		}

		c.Next()
	}

}

// Lê e valida os UUIDs no corpo da requisição, se houver conteúdo
func validateBodyUUIDs(c *gin.Context) error {
	// Verifica se há corpo na requisição
	if c.Request.ContentLength == 0 {
		return nil
	}

	// Faz o bind do JSON para um mapa
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}

	return validateUUIDsRecursive(body)
}

func validateUUIDsRecursive(data map[string]interface{}) error {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if strings.Contains(key, "id") {
				_, err := uuid.Parse(v)
				if err != nil {
					return errs.InvalidParam(key, err)
				}
			}
		case []any:
			for _, item := range v {
				if strItem, ok := item.(string); ok {
					_, err := uuid.Parse(strItem)
					if err != nil {
						return errs.InvalidParam(key, err)
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
