package services

import (
	"backend-go/internal/api/config"
	"backend-go/internal/api/errs"
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository"
	"backend-go/pkg/pagination"
	"backend-go/pkg/utils"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ParseDebt(debtReq dto.DebtRequest) (models.Debt, error) {
	purchaseDate, err := time.Parse("2006-01-02", debtReq.PurchaseDate)
	if err != nil {
		return models.Debt{}, errs.ParsingField("purchase_date", err)
	}

	dueDate, err := time.Parse("2006-01-02", debtReq.DueDate)
	if err != nil {
		return models.Debt{}, errs.ParsingField("due_date", err)
	}

	amount, err := strconv.ParseFloat(debtReq.Amount, 64)
	if err != nil {
		return models.Debt{}, errs.ParsingField("amount", err)
	}

	var categoryID *uuid.UUID
	categoryID = nil
	category := categorizeTransaction(debtReq.Title)
	if category != nil {
		categoryID, err = repository.GetCategoryIDByName(category)
		if errors.Is(err, errs.ErrNoRows) {
			return models.Debt{}, errs.NotFound("category", *category)
		}

	}
	if err != nil {
		return models.Debt{}, errs.UnknownWithContext("buscar categoria", err)
	}

	invoiceID, err := utils.ToUUIDPointer(debtReq.InvoiceID)
	if err != nil {
		return models.Debt{}, errs.ParsingField("invoice_id", err)
	}

	return models.Debt{
		InvoiceID:    invoiceID,
		Title:        debtReq.Title,
		Amount:       amount,
		PurchaseDate: purchaseDate,
		DueDate:      dueDate,
		CategoryID:   categoryID,
	}, nil
}

func categorizeTransaction(name string) *string {
	if category, exists := config.CategoryMap[name]; exists {
		return &category
	}
	return nil
}

func CreateDebt(debt models.Debt) (models.Debt, error) {
	return repository.InsertDebt(debt)
}

func UpdateDebt(debt models.Debt) (models.Debt, error) {
	return repository.UpdateDebt(debt)
}

func ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, int, error) {

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
		return nil, 0, err
	}

	debts, err := repository.ListDebts(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := repository.CountDebts()
	if err != nil {
		return nil, 0, err
	}

	return debts, total, nil

}

func GetDebtByID(id uuid.UUID) (*models.Debt, error) {
	return repository.GetDebtByID(id)
}

func DeleteDebtByID(id uuid.UUID) error {
	return repository.DeleteDebtByID(id)
}
