package interfaces

import (
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	// Debt
	// GetDebtByID(id uuid.UUID) (*models.Debt, error)
	// DeleteDebtByID(id uuid.UUID) error
	// InsertDebt(debt models.Debt) (models.Debt, error)
	// UpdateDebt(debt models.Debt) (models.Debt, error)
	// ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtsResponse, error)
	// CountDebts(flt dto.DebtFilters, pgn *pagination.Pagination) (int, error)
	// Invoice
	// GetInvoiceByID(id uuid.UUID) (*models.Invoice, error)
	// DeleteInvoiceByID(id uuid.UUID) error
	// InsertInvoice(input models.Invoice) (models.Invoice, error)
	// UpdateInvoice(input models.Invoice) (models.Invoice, error)
	// ListInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error)
	// CountInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error)
	// Category
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	InsertCategory(ctx context.Context, input models.Category) (*models.Category, error)
	UpdateCategory(ctx context.Context, input models.Category) (*models.Category, error)
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoriesResponse, error)
	CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error)
	// PaymentStatus
	// GetPaymentStatusByID(id uuid.UUID) (*models.PaymentStatus, error)
	// GetPaymentStatusIDByName(name *string) (*uuid.UUID, error)
	// DeletePaymentStatusByID(id uuid.UUID) error
	// InsertPaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error)
	// UpdatePaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error)
	// ListPaymentStatus(pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, error)
	// CountPaymentStatus(pgn *pagination.Pagination) (int, error)
}
