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
	// TODO: ele é obrigatorio no banco, ver depois como lidar com isso e o seu hook
	StatusID     *uuid.UUID  `json:"status_id"`
}

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

type Invoice struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	IssueDate time.Time `json:"issue_date"`
	DueDate   time.Time `json:"due_date"`
	StatusID  uuid.UUID `json:"status_id"`
}

type PaymentStatus struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}
