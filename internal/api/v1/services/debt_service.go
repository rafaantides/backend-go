package services

import (
	"backend-go/internal/api/config"
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"
	"backend-go/pkg/utils"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type DebtService struct {
	DB interfaces.Database
	// MQ interfaces.MessageQueue
}

func NewDebtService(db interfaces.Database) *DebtService {
	return &DebtService{DB: db}
}

func (s *DebtService) ParseDebt(debtReq dto.DebtRequest) (models.Debt, error) {
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
		categoryID, err = s.DB.GetCategoryIDByName(category)
		if errors.Is(err, errs.ErrNotFound) {
			return models.Debt{}, errs.ResorceNotFound("category", *category)
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

func (s *DebtService) CreateDebt(debt models.Debt) (models.Debt, error) {
	return s.DB.InsertDebt(debt)
}

func (s *DebtService) UpdateDebt(debt models.Debt) (models.Debt, error) {
	return s.DB.UpdateDebt(debt)
}

func (s *DebtService) ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtsResponse, int, error) {
	debts, err := s.DB.ListDebts(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountDebts(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return debts, total, nil

}

func (s *DebtService) GetDebtByID(id uuid.UUID) (*models.Debt, error) {
	return s.DB.GetDebtByID(id)
}

func (s *DebtService) DeleteDebtByID(id uuid.UUID) error {
	return s.DB.DeleteDebtByID(id)
}

func categorizeTransaction(name string) *string {
	if category, exists := config.CategoryMap[name]; exists {
		return &category
	}
	return nil
}
