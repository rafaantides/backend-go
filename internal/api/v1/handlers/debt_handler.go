package handlers

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/services"
	"backend-go/pkg/pagination"
	"backend-go/pkg/utils"
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
// @Failure 400 {object} errs.ErrorResponse "Requisição inválida"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts [post]
func CreateDebtHandler(c *gin.Context) {
	var debtReq dto.DebtRequest

	if err := c.ShouldBindJSON(&debtReq); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	debt, err := services.ParseDebt(debtReq)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	createdDebt, err := services.CreateDebt(debt)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
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
// @Failure 400 {object} errs.ErrorResponse "ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts/{id} [get]
func GetDebtByIDHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	debt, err := services.GetDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, debt)
}

// @Summary Listar todos os débitos
// @Description Retorna uma lista de débitos com paginação e filtros opcionais
// @Tags Débitos
// @Accept json
// @Produce json
// @Param title query string false "Filtrar por título do débito"
// @Param category_id query string false "Filtrar por ID da categoria (UUID)"
// @Param status_id query string false "Filtrar por ID do status (UUID)"
// @Param min_amount query number false "Valor mínimo do débito"
// @Param max_amount query number false "Valor máximo do débito"
// @Param start_date query string false "Filtrar por data de início (YYYY-MM-DD)"
// @Param end_date query string false "Filtrar por data de término (YYYY-MM-DD)"
// @Param invoice_id query string false "Filtrar por ID da fatura (UUID)"
// @Param page query int false "Número da página"
// @Param page_size query int false "Tamanho da página"
// @Param order_by query string false "Ordenação dos resultados (ex: amount, due_date)"
// @Success 200 {array} dto.DebtResponse "Lista de débitos"
// @Failure 400 {object} errs.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts [get]
func ListDebtsHandler(c *gin.Context) {
	var flt dto.DebtFilters
	if err := c.ShouldBindQuery(&flt); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	pgn, err := pagination.NewPagination(c)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	validColumns := map[string]bool{
		"id":            true,
		"invoice_id":    true,
		"title":         true,
		"category_id":   true,
		"amount":        true,
		"purchase_date": true,
		"due_date":      true,
		"status_id":     true,
		"created_at":    true,
		"updated_at":    true,
	}

	if err := pgn.ValidateOrderBy("purchase_date", validColumns); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	debts, total, err := services.ListDebts(flt, pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

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
// @Failure 400 {object} errs.ErrorResponse "Requisição inválida ou ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts/{id} [put]
func UpdateDebtHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	_, err = services.GetDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	var debtReq dto.DebtRequest
	if err := c.ShouldBindJSON(&debtReq); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	debt, err := services.ParseDebt(debtReq)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	updateDebt, err := services.UpdateDebt(debt)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
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
// @Failure 400 {object} errs.ErrorResponse "ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts/{id} [delete]
func DeleteDebtHandler(c *gin.Context) {
	debtID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || debtID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = services.DeleteDebtByID(*debtID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}

		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
