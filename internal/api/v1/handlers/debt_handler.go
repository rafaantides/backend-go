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

type DebtHandler struct {
	Service *services.DebtService
}

func NewDebtHandler(service *services.DebtService) *DebtHandler {
	return &DebtHandler{Service: service}
}

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
func (h *DebtHandler) CreateDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.DebtRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := h.Service.ParseDebt(ctx, req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	newDebt, err := h.Service.CreateDebt(ctx, input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, newDebt)
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
func (h *DebtHandler) GetDebtByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := h.Service.GetDebtByID(ctx, *id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
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
// @Success 200 {array} dto.DebtsResponse "Lista de débitos"
// @Failure 400 {object} errs.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /debts [get]
func (h *DebtHandler) ListDebtsHandler(c *gin.Context) {
	ctx := c.Request.Context()
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

	response, total, err := h.Service.ListDebts(ctx, flt, pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
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
func (h *DebtHandler) UpdateDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	_, err = h.Service.GetDebtByID(ctx, *id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	var req dto.DebtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := h.Service.ParseDebt(ctx, req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := h.Service.UpdateDebt(ctx, input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
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
func (h *DebtHandler) DeleteDebtHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = h.Service.DeleteDebtByID(ctx, *id)
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
