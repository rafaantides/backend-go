package handlers

import (
	"api-go/internal/models"
	"api-go/internal/services"
	"api-go/pkg/utils"
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

func GetAllDebtsHandler(c *gin.Context) {

	var filters models.DebtFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Parâmetros inválidos",
			Details: err.Error(),
		})
		return
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 10
	}

	debts, total, err := services.GetAllDebts(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Erro ao buscar dívidas",
			Details: err.Error(),
		})
		return
	}

	utils.SetPaginationHeaders(c, filters.Page, filters.PageSize, total)

	c.JSON(http.StatusOK, debts)
}
