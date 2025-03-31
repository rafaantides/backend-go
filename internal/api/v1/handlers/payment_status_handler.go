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

func CreatePaymentStatusHandler(c *gin.Context) {
	var req dto.PaymentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	input, err := services.ParsePaymentStatus(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.CreatePaymentStatus(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}

func GetPaymentStatusByIDHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.GetPaymentStatusByID(*id)
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

func ListPaymentStatussHandler(c *gin.Context) {
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

	response, total, err := services.ListCategories(pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, response)
}

func UpdatePaymentStatusHandler(c *gin.Context) {
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

	input, err := services.ParsePaymentStatus(req)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	data, err := services.UpdatePaymentStatus(input)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

func DeletePaymentStatusHandler(c *gin.Context) {
	id, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || id == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = services.DeletePaymentStatusByID(*id)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
