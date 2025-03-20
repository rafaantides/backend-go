package models

import (
	"time"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Debt struct {
	ID           uuid.UUID  `json:"id"`
	InvoiceID    *uuid.UUID `json:"invoice_id,omitempty"`
	Title        string     `json:"title"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      time.Time  `json:"due_date,omitempty"`
	StatusID     uuid.UUID  `json:"status_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type DebtRequest struct {
	InvoiceID    string `json:"invoice_id"`
	PurchaseDate string `json:"purchase_date"`
	Title        string `json:"title"`
	Amount       string `json:"amount"`
}

type DebtResponse struct {
	ID           uuid.UUID `json:"id"`
	InvoiceTitle *string   `json:"invoice_title,omitempty"`
	Title        string    `json:"title"`
	Amount       float64   `json:"amount"`
	PurchaseDate string    `json:"purchase_date"`
	DueDate      *string   `json:"due_date"`
	Category     *string   `json:"category,omitempty"`
	StatusID     uuid.UUID `json:"status_id"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type DebtFilters struct {
	Title      string     `form:"title"`
	CategoryID *uuid.UUID `form:"category_id"`
	StatusID   *uuid.UUID `form:"status_id"`
	MinAmount  *float64   `form:"min_amount"`
	MaxAmount  *float64   `form:"max_amount"`
	StartDate  string     `form:"start_date"`
	EndDate    string     `form:"end_date"`
	InvoiceID  *uuid.UUID `form:"invoice_id"`
	Page       int        `form:"page" binding:"required"`
	PageSize   int        `form:"page_size" binding:"required"`
}
