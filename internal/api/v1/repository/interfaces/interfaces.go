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
	GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error)
	DeleteDebtByID(ctx context.Context, id uuid.UUID) error
	InsertDebt(ctx context.Context, input models.Debt) (*dto.DebtResponse, error)
	UpdateDebt(ctx context.Context, input models.Debt) (*dto.DebtResponse, error)
	ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, error)
	CountDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) (int, error)
	// Invoice
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceResponse, error)
	DeleteInvoiceByID(ctx context.Context, id uuid.UUID) error
	InsertInvoice(ctx context.Context, input models.Invoice) (*dto.InvoiceResponse, error)
	UpdateInvoice(ctx context.Context, input models.Invoice) (*dto.InvoiceResponse, error)
	ListInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error)
	CountInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error)
	// Category
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	InsertCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error)
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, error)
	CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error)
	// PaymentStatus
	GetPaymentStatusByID(ctx context.Context, id uuid.UUID) (*dto.PaymentStatusResponse, error)
	GetPaymentStatusIDByName(ctx context.Context, name *string) (*uuid.UUID, error)
	DeletePaymentStatusByID(ctx context.Context, id uuid.UUID) error
	InsertPaymentStatus(ctx context.Context, input models.PaymentStatus) (*dto.PaymentStatusResponse, error)
	UpdatePaymentStatus(ctx context.Context, input models.PaymentStatus) (*dto.PaymentStatusResponse, error)
	ListPaymentStatus(ctx context.Context, pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, error)
	CountPaymentStatus(ctx context.Context, pgn *pagination.Pagination) (int, error)
}
