package models

import (
	"time"

	"github.com/google/uuid"
)

type Debt struct {
	ID           uuid.UUID  `json:"id"`
	InvoiceID    *uuid.UUID `json:"invoice_id"`
	Title        string     `json:"title"`
	CategoryID   *uuid.UUID `json:"category_id"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      time.Time  `json:"due_date"`
	StatusID     uuid.UUID  `json:"status_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Invoice struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	IssueDate time.Time `json:"issue_date"`
	DueDate   time.Time `json:"due_date"`
	StatusID  uuid.UUID `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaymentStatus struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
