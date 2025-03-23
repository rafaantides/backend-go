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

// @Summary Criar uma nova fatura
// @Description Cria uma nova fatura com os dados fornecidos no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param invoice body dto.InvoiceRequest true "Dados da fatura"
// @Success 201 {object} models.Invoice
// @Failure 400 {object} dto.ErrorResponse "Requisição inválida"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/invoices [post]
func CreateInvoiceHandler(c *gin.Context) {
	var invoiceReq dto.InvoiceRequest

	if err := c.ShouldBindJSON(&invoiceReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	invoice, err := services.ParseInvoice(invoiceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Dados inválidos",
			Details: err.Error(),
		})
		return
	}

	createdInvoice, err := services.CreateInvoice(invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao salvar a fatura",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdInvoice)
}

// @Summary Buscar fatura por ID
// @Description Retorna uma fatura pelo ID fornecido na URL
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} dto.ErrorResponse "ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/invoices/{id} [get]
func GetInvoiceByIDHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID da fatura inválido",
			Details: err.Error(),
		})
		return
	}

	invoice, err := services.GetInvoiceByID(*invoiceID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Fatura não encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao buscar a fatura",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invoice)
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
// @Failure 400 {object} dto.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/invoices [get]
func ListInvoicesHandler(c *gin.Context) {
	var filters dto.InvoiceFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Parâmetros inválidos",
			Details: err.Error(),
		})
		return
	}

	page := pagination.GetPage(filters.Page)
	pageSize := pagination.GetPageSize(filters.PageSize)

	invoices, total, err := services.ListInvoices(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao buscar faturas",
			Details: err.Error(),
		})
		return
	}

	pagination.SetPaginationHeaders(c, page, pageSize, total)

	c.JSON(http.StatusOK, invoices)
}

// @Summary Atualizar uma fatura
// @Description Atualiza uma fatura existente com os novos dados fornecidos no corpo da requisição
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Param invoice body dto.InvoiceRequest true "Dados da fatura"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} dto.ErrorResponse "Requisição inválida ou ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/invoices/{id} [put]
func UpdateInvoiceHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID da fatura inválido",
			Details: err.Error(),
		})
		return
	}

	var invoiceReq dto.InvoiceRequest
	if err := c.ShouldBindJSON(&invoiceReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Requisição inválida",
			Details: err.Error(),
		})
		return
	}

	invoice, err := services.ParseInvoice(invoiceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Dados inválidos",
			Details: err.Error(),
		})
		return
	}

	updatedInvoice, err := services.UpdateInvoice(invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao atualizar a fatura",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedInvoice)
}

// @Summary Deletar uma fatura
// @Description Remove uma fatura pelo ID fornecido
// @Tags Faturas
// @Accept json
// @Produce json
// @Param id path string true "ID da fatura"
// @Success 204 "Registro deletado com sucesso"
// @Failure 400 {object} dto.ErrorResponse "ID inválido"
// @Failure 404 {object} dto.ErrorResponse "Registro não encontrado"
// @Failure 500 {object} dto.ErrorResponse "Erro interno"
// @Router /api/invoices/{id} [delete]
func DeleteInvoiceHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "ID da fatura inválido",
			Details: err.Error(),
		})
		return
	}

	err = services.DeleteInvoiceByID(*invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Erro ao deletar a fatura",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
