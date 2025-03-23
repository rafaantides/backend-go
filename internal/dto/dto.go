package dto

import "github.com/google/uuid"

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Debts
type DebtRequest struct {
	InvoiceID    string `json:"invoice_id"`
	PurchaseDate string `json:"purchase_date"`
	DueDate      string `json:"due_date"`
	Title        string `json:"title"`
	Amount       string `json:"amount"`
}

type DebtResponse struct {
	// ID único do débito
	ID uuid.UUID `json:"id"`
	// Título da fatura associada
	InvoiceTitle *string `json:"invoice_title"`
	// Título do débito
	Title string `json:"title"`
	// Valor do débito
	Amount float64 `json:"amount"`
	// Data da compra no formato YYYY-MM-DD
	PurchaseDate string `json:"purchase_date"`
	// Data de vencimento no formato YYYY-MM-DD
	DueDate *string `json:"due_date"`
	// ID da categoria do débito
	CategoryID *uuid.UUID `json:"category_id"`
	// Nome da categoria do débito
	Category *string `json:"category"`
	// ID do status do débito
	StatusID uuid.UUID `json:"status_id"`
	// Nome do status do débito
	Status string `json:"status"`
	// Data de criação do débito
	CreatedAt string `json:"created_at"`
	// Data da última atualização do débito
	UpdatedAt string `json:"updated_at"`
}

type DebtFilters struct {
	Title      *string    `form:"title"`
	CategoryID *uuid.UUID `form:"category_id"`
	StatusID   *uuid.UUID `form:"status_id"`
	MinAmount  *float64   `form:"min_amount"`
	MaxAmount  *float64   `form:"max_amount"`
	StartDate  *string    `form:"start_date"`
	EndDate    *string    `form:"end_date"`
	InvoiceID  *uuid.UUID `form:"invoice_id"`
	Page       *int       `form:"page"`
	PageSize   *int       `form:"page_size"`
	OrderBy    *string    `form:"order_by"`
}

// Invoices

type InvoiceRequest struct {
	Title     string `json:"title"`
	Amount    string `json:"amount"`
	IssueDate string `json:"issue_date"`
	DueDate   string `json:"due_date"`
}

type InvoiceResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	IssueDate string    `json:"issue_date"`
	DueDate   *string   `json:"due_date"`
	StatusID  uuid.UUID `json:"status_id"`
	Status    *string   `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type InvoiceFilters struct {
	Title     *string    `form:"title"`
	StatusID  *uuid.UUID `form:"status_id"`
	MinAmount *float64   `form:"min_amount"`
	MaxAmount *float64   `form:"max_amount"`
	StartDate *string    `form:"start_date"`
	EndDate   *string    `form:"end_date"`
	Page      *int       `form:"page"`
	PageSize  *int       `form:"page_size"`
	OrderBy   *string    `form:"order_by"`
}
