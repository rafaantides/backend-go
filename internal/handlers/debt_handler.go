package handlers

import (
	"api-go/internal/errs"
	"api-go/internal/models"
	"api-go/internal/services"
	"api-go/pkg/utils"
	"errors"
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

	c.JSON(http.StatusCreated, createdDebt)
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

func UpdateDebtHandler(c *gin.Context) {

	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "ID do débito inválido",
			Details: err.Error(),
		})
		return
	}

	debtDB, err := services.GetDebtByID(*debtID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Erro ao buscar o debito",
			Details: err.Error(),
		})
		return
	}

	if debtDB == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "Debito não encontrado",
		})
		return
	}

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

	updateDebt, err := services.UpdateDebt(debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Erro ao atualizar o debito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, updateDebt)
}

func DeleteDebtHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "ID do débito inválido",
			Details: err.Error(),
		})
		return
	}

	err = services.DeleteDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Message: "Débito não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Erro ao deletar o débito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
