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

func CreateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return repository.InsertInvoice(Invoice)
}

func UpdateInvoice(Invoice models.Invoice) (models.Invoice, error) {
	return repository.UpdateInvoice(Invoice)
}

func ListInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, int, error) {
	invoices, err := repository.ListInvoices(flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := repository.CountInvoices(flt, pgn)
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
