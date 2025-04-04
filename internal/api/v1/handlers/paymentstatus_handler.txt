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

type PaymentStatusHandler struct {
	Service *services.PaymentStatusService
}

func NewPaymentStatusHandler(service *services.PaymentStatusService) *PaymentStatusHandler {
	return &PaymentStatusHandler{Service: service}
}

func (h *PaymentStatusHandler) CreatePaymentStatusHandler(c *gin.Context) {
	var req dto.PaymentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := h.Service.ParsePaymentStatus(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := h.Service.CreatePaymentStatus(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

func (h *PaymentStatusHandler) GetPaymentStatusByIDHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := h.Service.GetPaymentStatusByID(*id)
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

func (h *PaymentStatusHandler) ListPaymentStatussHandler(c *gin.Context) {
	pgn, err := pagination.NewPagination(c)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	validColumns := map[string]bool{
		"id":          true,
		"name":        true,
		"description": true,
	}

	if err := pgn.ValidateOrderBy("name", validColumns); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	response, total, err := h.Service.ListPaymentStatus(pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

func (h *PaymentStatusHandler) UpdatePaymentStatusHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	var req dto.PaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := h.Service.ParsePaymentStatus(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := h.Service.UpdatePaymentStatus(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PaymentStatusHandler) DeletePaymentStatusHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = h.Service.DeletePaymentStatusByID(*id)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
