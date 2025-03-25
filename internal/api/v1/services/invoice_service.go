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

func ParseInvoice(invoiceReq dto.InvoiceRequest) (models.Invoice, error) {
	issueDate, err := time.Parse("2006-01-02", invoiceReq.IssueDate)
	if err != nil {
		return models.Invoice{}, errs.ErrDateParsing("issue_date")
	}

	dueDate, err := time.Parse("2006-01-02", invoiceReq.DueDate)
	if err != nil {
		return models.Invoice{}, errs.ErrDateParsing("due_date")
	}

	amount, err := strconv.ParseFloat(invoiceReq.Amount, 64)
	if err != nil {
		return models.Invoice{}, errs.ErrParsingField("amount", err)
	}

	return models.Invoice{
		Title:     invoiceReq.Title,
		Amount:    amount,
		IssueDate: issueDate,
		DueDate:   dueDate,
	}, nil

}

func CreateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return repository.InsertInvoice(Invoice)
}

func UpdateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return repository.UpdateInvoice(Invoice)
}

func ListInvoices(filters dto.InvoiceFilters, pagination *pagination.Pagination) ([]dto.InvoiceResponse, int, error) {

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

	if err := pagination.ValidateOrderBy("issue_date", validColumns); err != nil {
		return nil, 0, err
	}

	invoices, err := repository.ListInvoices(filters, pagination)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := repository.CountInvoices()
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func GetInvoiceByID(id uuid.UUID) (*models.Invoice, error) {
	return repository.GetInvoiceByID(id)
}

func DeleteInvoiceByID(id uuid.UUID) error {
	return repository.DeleteInvoiceByID(id)
}
