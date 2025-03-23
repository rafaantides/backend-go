package handlers

import (
	"api-go/internal/dto"
	"api-go/internal/errs"
	"api-go/internal/services"
	"api-go/pkg/pagination"
	"api-go/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Criar um novo débito
// @Description Cria um novo débito com os dados fornecidos no corpo da requisição
// @Tags Débitos
// @Accept json
// @Produce json
// @Param debt body dto.DebtRequest true "Dados do débito"
// @Success 201 {object} models.Debt
// @Failure 400 {object} dto.ErrorResponse "Requisição inválida"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/debts [post]
func CreateDebtHandler(c *gin.Context) {
	var debtReq dto.DebtRequest

	if err := c.ShouldBindJSON(&debtReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	debt, err := services.ParseDebt(debtReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Dados inválidos",
			Details: err.Error(),
		})
		return
	}

	createdDebt, err := services.CreateDebt(debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao salvar a débito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdDebt)
}

// @Summary Buscar débito por ID
// @Description Retorna um débito pelo ID fornecido na URL
// @Tags Débitos
// @Accept json
// @Produce json
// @Param id path string true "ID do débito"
// @Success 200 {object} models.Debt
// @Failure 400 {object} dto.ErrorResponse "ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/debts/{id} [get]
func GetDebtByIDHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID do débito inválido",
			Details: err.Error(),
		})
		return
	}

	debt, err := services.GetDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Débito não encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao buscar o débito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, debt)
}

// @Summary Listar todos os débitos
// @Description Retorna uma lista de débitos com paginação
// @Tags Débitos
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param pageSize query int false "Tamanho da página"
// @Success 200 {array} models.Debt
// @Failure 400 {object} dto.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/debts [get]
func ListDebtsHandler(c *gin.Context) {
	var filters dto.DebtFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Parâmetros inválidos",
			Details: err.Error(),
		})
		return
	}

	page := pagination.GetPage(filters.Page)
	pageSize := pagination.GetPageSize(filters.PageSize)

	debts, total, err := services.ListDebts(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao buscar os débitos",
			Details: err.Error(),
		})
		return
	}

	pagination.SetPaginationHeaders(c, page, pageSize, total)

	c.JSON(http.StatusOK, debts)
}

// @Summary Atualizar um débito
// @Description Atualiza um débito existente com os novos dados fornecidos no corpo da requisição
// @Tags Débitos
// @Accept json
// @Produce json
// @Param id path string true "ID do débito"
// @Param debt body dto.DebtRequest true "Dados do débito"
// @Success 200 {object} models.Debt
// @Failure 400 {object} dto.ErrorResponse "Requisição inválida ou ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/debts/{id} [put]
func UpdateDebtHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID do débito inválido",
			Details: err.Error(),
		})
		return
	}

	_, err = services.GetDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Débito não encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao buscar o débito",
			Details: err.Error(),
		})
		return
	}

	var debtReq dto.DebtRequest
	if err := c.ShouldBindJSON(&debtReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	debt, err := services.ParseDebt(debtReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Dados inválidos",
			Details: err.Error(),
		})
		return
	}

	updateDebt, err := services.UpdateDebt(debt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao atualizar o débito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updateDebt)
}

// @Summary Deletar um débito
// @Description Remove um débito pelo ID fornecido
// @Tags Débitos
// @Accept json
// @Produce json
// @Param id path string true "ID do débito"
// @Success 204 "Registro deletado com sucesso"
// @Failure 400 {object} dto.ErrorResponse "ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/debts/{id} [delete]
func DeleteDebtHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID do débito inválido",
			Details: err.Error(),
		})
		return
	}

	err = services.DeleteDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Débito não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao deletar o débito",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
