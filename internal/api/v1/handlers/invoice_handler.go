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

// @Summary Criar uma nova fatura
// @Description Cria uma nova fatura com os dados fornecidos no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param invoice body dto.InvoiceRequest true "Dados da fatura"
// @Success 201 {object} models.Invoice
// @Failure 400 {object} errs.ErrorResponse "Requisição inválida"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /invoices [post]
func CreateInvoiceHandler(c *gin.Context) {
	var req dto.InvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := services.ParseInvoice(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.CreateInvoice(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

// @Summary Buscar fatura por ID
// @Description Retorna uma fatura pelo ID fornecido na URL
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} errs.ErrorResponse "ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /invoices/{id} [get]
func GetInvoiceByIDHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.GetInvoiceByID(*id)
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

// @Summary Listar faturas
// @Description Retorna uma lista de faturas com filtros opcionais
// @Tags Faturas
// @Accept json
// @Produce json
// @Param title query string false "Título da fatura"
// @Param status_id query string false "ID do status da fatura (UUID)"
// @Param min_amount query number false "Valor mínimo da fatura"
// @Param max_amount query number false "Valor máximo da fatura"
// @Param start_date query string false "Data inicial para filtrar (YYYY-MM-DD)"
// @Param end_date query string false "Data final para filtrar (YYYY-MM-DD)"
// @Param page query integer false "Número da página"
// @Param page_size query integer false "Tamanho da página"
// @Param order_by query string false "Campo de ordenação (ex: title, amount)"
// @Success 200 {array} dto.InvoiceResponse "Lista de faturas"
// @Failure 400 {object} errs.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /invoices [get]
func ListInvoicesHandler(c *gin.Context) {
	var flt dto.InvoiceFilters
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
		"id":         true,
		"title":      true,
		"amount":     true,
		"issue_date": true,
		"due_date":   true,
		"status_id":  true,
		"created_at": true,
		"updated_at": true,
	}

	if err := pgn.ValidateOrderBy("issue_date", validColumns); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	response, total, err := services.ListInvoices(flt, pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

// @Summary Atualizar uma fatura
// @Description Atualiza uma fatura existente com os novos dados fornecidos no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Param invoice body dto.InvoiceRequest true "Dados da fatura"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} errs.ErrorResponse "Requisição inválida ou ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /invoices/{id} [put]
func UpdateInvoiceHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	var req dto.InvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := services.ParseInvoice(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.UpdateInvoice(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary Deletar uma fatura
// @Description Remove uma fatura pelo ID fornecido
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 204 "Registro deletado com sucesso"
// @Failure 400 {object} errs.ErrorResponse "ID inválido"
// @Failure 404 {object} errs.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} errs.ErrorResponse "Erro interno"
// @Router /invoices/{id} [delete]
func DeleteInvoiceHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = services.DeleteInvoiceByID(*id)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
