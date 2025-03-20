package handlers

import (
	"api-go/internal/models"
	"api-go/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDebtHandler(c *gin.Context) {
	var debtReq models.DebtRequest

	if err := c.ShouldBindJSON(&debtReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	debt, err := services.ParseDebt(debtReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Dados inválidos",
			Details: err.Error(),
		})
		return
	}

	createdDebt, err := services.CreateDebt(debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Erro ao salvar a debito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Débito cadastrado com sucesso",
		"debt":    createdDebt,
	})
}
