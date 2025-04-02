package interfaces

import (
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	// Debt
	GetDebtByID(id uuid.UUID) (*models.Debt, error)
	DeleteDebtByID(id uuid.UUID) error
	InsertDebt(debt models.Debt) (models.Debt, error)
	UpdateDebt(debt models.Debt) (models.Debt, error)
	ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtsResponse, error)
	CountDebts(flt dto.DebtFilters, pgn *pagination.Pagination) (int, error)
	// Invoice
	GetInvoiceByID(id uuid.UUID) (*models.Invoice, error)
	DeleteInvoiceByID(id uuid.UUID) error
	InsertInvoice(input models.Invoice) (models.Invoice, error)
	UpdateInvoice(input models.Invoice) (models.Invoice, error)
	ListInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error)
	CountInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error)
	// Category
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	GetCategoryIDByName(name *string) (*uuid.UUID, error)
	DeleteCategoryByID(id uuid.UUID) error
	InsertCategory(input models.Category) (models.Category, error)
	UpdateCategory(input models.Category) (models.Category, error)
	ListCategories(pgn *pagination.Pagination) ([]dto.CategoriesResponse, error)
	CountCategories(pgn *pagination.Pagination) (int, error)
	// PaymentStatus
	GetPaymentStatusByID(id uuid.UUID) (*models.PaymentStatus, error)
	GetPaymentStatusIDByName(name *string) (*uuid.UUID, error)
	DeletePaymentStatusByID(id uuid.UUID) error
	InsertPaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error)
	UpdatePaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error)
	ListPaymentStatus(pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, error)
	CountPaymentStatus(pgn *pagination.Pagination) (int, error)
}
