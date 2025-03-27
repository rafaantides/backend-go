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

func CreateCategoryHandler(c *gin.Context) {
	var invoiceReq dto.CategoryRequest

	if err := c.ShouldBindJSON(&invoiceReq); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	invoice, err := services.ParseCategory(invoiceReq)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	newCategory, err := services.CreateCategory(invoice)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

func GetCategoryByIDHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	invoice, err := services.GetCategoryByID(*invoiceID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Error(errs.NewAPIError(http.StatusNotFound, err))
			return
		}
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func ListCategorysHandler(c *gin.Context) {
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

	if err := pgn.ValidateOrderBy("issue_date", validColumns); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	invoices, total, err := services.ListCategories(pgn)

	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	pgn.SetPaginationHeaders(c, total)

	c.JSON(http.StatusOK, invoices)
}

func UpdateCategoryHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	var invoiceReq dto.CategoryRequest
	if err := c.ShouldBindJSON(&invoiceReq); err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	invoice, err := services.ParseCategory(invoiceReq)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	updatedCategory, err := services.UpdateCategory(invoice)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func DeleteCategoryHandler(c *gin.Context) {
	invoiceID, err := utils.ToUUIDPointer(c.Param("id"))
	if err != nil || invoiceID == nil {
		c.Error(errs.NewAPIError(http.StatusBadRequest, err))
		return
	}

	err = services.DeleteCategoryByID(*invoiceID)
	if err != nil {
		c.Error(errs.NewAPIError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
