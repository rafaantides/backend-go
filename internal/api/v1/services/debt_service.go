package services

import (
	"backend-go/internal/api/config"
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	queue "backend-go/internal/api/v1/queue/interfaces"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"context"
	"log"

	"backend-go/pkg/pagination"
	"backend-go/pkg/utils"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type DebtService struct {
	DB repository.Database
	MQ queue.MessageQueue
}

func NewDebtService(db repository.Database, mq queue.MessageQueue) *DebtService {
	return &DebtService{DB: db, MQ: mq}
}

func (s *DebtService) ParseDebt(ctx context.Context, debtReq dto.DebtRequest) (models.Debt, error) {
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
		categoryID, err = s.DB.GetCategoryIDByName(ctx, category)
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

func (s *DebtService) CreateDebt(ctx context.Context, debt models.Debt) (*dto.DebtResponse, error) {
	return s.DB.InsertDebt(ctx, debt)
}

func (s *DebtService) UpdateDebt(ctx context.Context, debt models.Debt) (*dto.DebtResponse, error) {
	return s.DB.UpdateDebt(ctx, debt)
}

func (s *DebtService) ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, int, error) {
	debts, err := s.DB.ListDebts(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountDebts(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return debts, total, nil
}

func (s *DebtService) GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error) {
	return s.DB.GetDebtByID(ctx, id)
}

func (s *DebtService) DeleteDebtByID(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeleteDebtByID(ctx, id)
}

func categorizeTransaction(name string) *string {
	if category, exists := config.CategoryMap[name]; exists {
		return &category
	}
	return nil
}
