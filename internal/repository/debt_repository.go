package repository

import (
	"api-go/internal/models"
	"fmt"
)

func InsertDebt(debt models.Debt) (models.Debt, error) {
	query := `INSERT INTO debts (invoice_id, title, category_id, amount, purchase_date, due_date, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
			  RETURNING id, invoice_id, title, category_id, amount, purchase_date, due_date, status_id, created_at, updated_at`

	var newDebt models.Debt
	err := DB.QueryRow(query, debt.InvoiceID, debt.Title, debt.CategoryID, debt.Amount, debt.PurchaseDate, debt.DueDate).
		Scan(&newDebt.ID, &newDebt.InvoiceID, &newDebt.Title, &newDebt.CategoryID, &newDebt.Amount, &newDebt.PurchaseDate, &newDebt.DueDate, &newDebt.StatusID, &newDebt.CreatedAt, &newDebt.UpdatedAt)
	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to insert debt: %w", err)
	}
	return newDebt, nil
}
