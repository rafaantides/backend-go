package services

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository"
	"backend-go/pkg/pagination"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type InvoiceService struct {
	DB *repository.Database
}

func NewInvoiceService(db *repository.Database) *InvoiceService {
	return &InvoiceService{DB: db}
}

func (s *InvoiceService) ParseInvoice(invoiceReq dto.InvoiceRequest) (models.Invoice, error) {
	issueDate, err := time.Parse("2006-01-02", invoiceReq.IssueDate)
	if err != nil {
		return models.Invoice{}, errs.DateParsing("issue_date")
	}

	dueDate, err := time.Parse("2006-01-02", invoiceReq.DueDate)
	if err != nil {
		return models.Invoice{}, errs.DateParsing("due_date")
	}

	amount, err := strconv.ParseFloat(invoiceReq.Amount, 64)
	if err != nil {
		return models.Invoice{}, errs.ParsingField("amount", err)
	}

	return models.Invoice{
		Title:     invoiceReq.Title,
		Amount:    amount,
		IssueDate: issueDate,
		DueDate:   dueDate,
	}, nil

}

func (s *InvoiceService) CreateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return s.DB.InsertInvoice(Invoice)
}

func (s *InvoiceService) UpdateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return s.DB.UpdateInvoice(Invoice)
}

func (s *InvoiceService) ListInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, int, error) {
	invoices, err := s.DB.ListInvoices(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountInvoices(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (s *InvoiceService) GetInvoiceByID(id uuid.UUID) (*models.Invoice, error) {
	return s.DB.GetInvoiceByID(id)
}

func (s *InvoiceService) DeleteInvoiceByID(id uuid.UUID) error {
	return s.DB.DeleteInvoiceByID(id)
}
