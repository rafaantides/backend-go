package repository

import (
	"api-go/internal/models"
	"database/sql"
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

func GetAllDebts(filters models.DebtFilters) ([]models.DebtResponse, int, error) {
	query := `
        SELECT 
            d.id, d.title, d.amount, d.purchase_date, d.due_date, d.status_id, 
            d.created_at, d.updated_at, c.name AS category, i.title AS invoice_title
        FROM debts d
        LEFT JOIN categories c ON d.category_id = c.id
        LEFT JOIN invoices i ON d.invoice_id = i.id
    `

	var conditions []string
	var args []any
	argIndex := 1

	if filters.Title != "" {
		conditions = append(conditions, fmt.Sprintf("d.title ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Title+"%")
		argIndex++
	}
	if filters.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("d.category_id = $%d", argIndex))
		args = append(args, *filters.CategoryID)
		argIndex++
	}
	if filters.StatusID != nil {
		conditions = append(conditions, fmt.Sprintf("d.status_id = $%d", argIndex))
		args = append(args, *filters.StatusID)
		argIndex++
	}
	if filters.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}
	if filters.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}
	if filters.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date >= $%d", argIndex))
		args = append(args, filters.StartDate)
		argIndex++
	}
	if filters.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date <= $%d", argIndex))
		args = append(args, filters.EndDate)
		argIndex++
	}
	if filters.InvoiceID != nil {
		conditions = append(conditions, fmt.Sprintf("d.invoice_id = $%d", argIndex))
		args = append(args, *filters.InvoiceID)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY d.purchase_date DESC LIMIT $%d OFFSET $%d"
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows, err := DB.Query(fmt.Sprintf(query, argIndex, argIndex+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var debts []models.DebtResponse
	for rows.Next() {
		var debt models.DebtResponse
		var dueDate sql.NullString
		var category, invoiceTitle sql.NullString

		err := rows.Scan(
			&debt.ID, &debt.Title, &debt.Amount, &debt.PurchaseDate, &dueDate,
			&debt.StatusID, &debt.CreatedAt, &debt.UpdatedAt, &category, &invoiceTitle,
		)
		if err != nil {
			return nil, 0, err
		}

		if dueDate.Valid {
			debt.DueDate = &dueDate.String
		}
		if category.Valid {
			debt.Category = &category.String
		}
		if invoiceTitle.Valid {
			debt.InvoiceTitle = &invoiceTitle.String
		}

		debts = append(debts, debt)
	}

	countQuery := "SELECT COUNT(*) FROM debts"
	var total int
	err = DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return debts, total, nil
}
