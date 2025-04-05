package dto

import "github.com/google/uuid"

// Debts
type DebtRequest struct {
	InvoiceID    string `json:"invoice_id"`
	PurchaseDate string `json:"purchase_date"`
	DueDate      string `json:"due_date"`
	Title        string `json:"title"`
	Amount       string `json:"amount"`
}

type DebtsResponse struct {
	// ID único do débito
	ID uuid.UUID `json:"id"`
	// Título da fatura associada
	// TODO: colocar invoiceID
	InvoiceTitle *string `json:"invoice_title"`
	// Título do débito
	Title string `json:"title"`
	// Valor do débito
	Amount float64 `json:"amount"`
	// Data da compra no formato YYYY-MM-DD
	PurchaseDate string `json:"purchase_date"`
	// Data de vencimento no formato YYYY-MM-DD
	DueDate *string `json:"due_date"`
	// ID da categoria
	CategoryID *uuid.UUID `json:"category_id"`
	// Nome da categoria
	Category *string `json:"category"`
	// ID do status
	// TODO: ver como manter obrigatorio
	StatusID *uuid.UUID `json:"status_id"`
	// Nome do status
	Status *string `json:"status"`
	// Data de criação do débito
	CreatedAt string `json:"created_at"`
	// Data da última atualização do débito
	UpdatedAt string `json:"updated_at"`
}

type DebtFilters struct {
	// REDO: o ShouldBindQuery n esta reconhecendo o *[]uuid.UUID
	CategoryID *[]string `form:"category_id"`
	StatusID   *[]string `form:"status_id"`
	InvoiceID  *[]string `form:"invoice_id"`
	MinAmount  *float64  `form:"min_amount"`
	MaxAmount  *float64  `form:"max_amount"`
	StartDate  *string   `form:"start_date"`
	EndDate    *string   `form:"end_date"`
}

// Invoices
type InvoiceRequest struct {
	Title     string `json:"title"`
	Amount    string `json:"amount"`
	IssueDate string `json:"issue_date"`
	DueDate   string `json:"due_date"`
}

type InvoiceResponse struct {
	// ID único da fatura
	ID uuid.UUID `json:"id"`
	// Título da fatura
	Title string `json:"title"`
	// Valor da fatura
	Amount float64 `json:"amount"`
	// Data de emissão no formato YYYY-MM-DD
	IssueDate string `json:"issue_date"`
	// Data de vencimento no formato YYYY-MM-DD
	DueDate *string `json:"due_date"`
	// ID do status
	StatusID uuid.UUID `json:"status_id"`
	// Nome do status
	Status string `json:"status"`
	// Data de criação da fatura
	CreatedAt string `json:"created_at"`
	// Data da última atualização da fatura
	UpdatedAt string `json:"updated_at"`
}

type InvoiceFilters struct {
	StatusID  *[]string `form:"status_id"`
	MinAmount *float64  `form:"min_amount"`
	MaxAmount *float64  `form:"max_amount"`
	StartDate *string   `form:"start_date"`
	EndDate   *string   `form:"end_date"`
}

// Category
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoriesResponse struct {
	// ID único da categoria
	ID uuid.UUID `json:"id"`
	// Nome da categoria
	Name string `json:"name"`
	// Descrião da categoria
	Description *string `json:"description"`
}

// PaymentStatus
type PaymentStatusRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PaymentStatusResponse struct {
	// ID único do status
	ID uuid.UUID `json:"id"`
	// Nome do status
	Name string `json:"name"`
	// Descrião do status
	Description float64 `json:"description"`
}
