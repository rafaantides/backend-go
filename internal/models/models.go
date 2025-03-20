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
