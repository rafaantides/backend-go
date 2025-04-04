package services

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type InvoiceService struct {
	DB repository.Database
}

func NewInvoiceService(db repository.Database) *InvoiceService {
	return &InvoiceService{DB: db}
}

func (s *InvoiceService) ParseInvoice(req dto.InvoiceRequest) (models.Invoice, error) {
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return models.Invoice{}, errs.DateParsing("issue_date")
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return models.Invoice{}, errs.DateParsing("due_date")
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return models.Invoice{}, errs.ParsingField("amount", err)
	}

	return models.Invoice{
		Title:     req.Title,
		Amount:    amount,
		IssueDate: issueDate,
		DueDate:   dueDate,
	}, nil

}

func (s *InvoiceService) CreateInvoice(ctx context.Context, input models.Invoice) (*models.Invoice, error) {
	return s.DB.InsertInvoice(ctx, input)
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, input models.Invoice) (*models.Invoice, error) {
	return s.DB.UpdateInvoice(ctx, input)
}

func (s *InvoiceService) ListInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, int, error) {
	invoices, err := s.DB.ListInvoices(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountInvoices(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (s *InvoiceService) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	return s.DB.GetInvoiceByID(ctx, id)
}

func (s *InvoiceService) DeleteInvoiceByID(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeleteInvoiceByID(ctx, id)
}
