package repository

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func GetDebtByID(id uuid.UUID) (*models.Debt, error) {
	var data models.Debt

	query := `SELECT * FROM debts WHERE id = $1`
	row := DB.QueryRow(query, id)

	err := row.Scan(&data.ID, &data.Title, &data.CategoryID, &data.Amount, &data.PurchaseDate, &data.DueDate, &data.StatusID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNoRows
		}
		return nil, err
	}

	return &data, nil
}

func DeleteDebtByID(id uuid.UUID) error {
	query := `DELETE FROM debts WHERE id = $1`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrNoRows
	}

	return nil
}

func InsertDebt(debt models.Debt) (models.Debt, error) {
	query := `INSERT INTO debts (invoice_id, title, category_id, amount, purchase_date, due_date)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING *`

	var newData models.Debt
	err := DB.QueryRow(query, debt.InvoiceID, debt.Title, debt.CategoryID, debt.Amount, debt.PurchaseDate, debt.DueDate).
		Scan(&newData.ID, &newData.InvoiceID, &newData.Title, &newData.CategoryID, &newData.Amount, &newData.PurchaseDate, &newData.DueDate, &newData.StatusID, &newData.CreatedAt, &newData.UpdatedAt)
	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to insert debt: %w", err)
	}
	return newData, nil
}

func UpdateDebt(debt models.Debt) (models.Debt, error) {
	query := `
		UPDATE debts 
		SET title = $1, amount = $2, purchase_date = $3, category_id = $4 , status_id = $5
		WHERE id = $6
		RETURNING *
	`
	var updatedData models.Debt
	err := DB.QueryRow(query, debt.Title, debt.Amount, debt.PurchaseDate, debt.CategoryID, debt.StatusID, debt.ID).
		Scan(&updatedData.ID, &updatedData.Title, &updatedData.Amount, &updatedData.PurchaseDate, &updatedData.CategoryID, &updatedData.StatusID, &updatedData.CreatedAt, &updatedData.UpdatedAt)

	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to update debt: %w", err)
	}

	return updatedData, nil
}

func ListDebts(filters dto.DebtFilters, page int, pageSize int, orderBy string) ([]dto.DebtResponse, int, error) {
	query := `
        SELECT
            d.id,
			d.title,
			d.amount,
			d.purchase_date,
			d.due_date,
			d.category_id,
			d.status_id,
            d.created_at,
			d.updated_at,
			c.name AS category,
			i.title AS invoice_title,
			s.name AS status
        FROM debts d
        LEFT JOIN categories c ON d.category_id = c.id
		LEFT JOIN payment_status s ON d.status_id = s.id
        LEFT JOIN invoices i ON d.invoice_id = i.id
    `

	var conditions []string
	var args []any
	argIndex := 1

	if filters.Title != nil {
		conditions = append(conditions, fmt.Sprintf("d.title ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Title+"%")
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
	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}
	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date <= $%d", argIndex))
		args = append(args, *filters.EndDate)
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

	query += fmt.Sprintf(" ORDER BY d.%s DESC", orderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var debts []dto.DebtResponse
	debts = make([]dto.DebtResponse, 0)
	for rows.Next() {
		var debt dto.DebtResponse

		err := rows.Scan(
			&debt.ID, &debt.Title, &debt.Amount, &debt.PurchaseDate, &debt.DueDate,
			&debt.CategoryID, &debt.StatusID, &debt.CreatedAt, &debt.UpdatedAt,
			&debt.Category, &debt.InvoiceTitle, &debt.Status,
		)
		if err != nil {
			return nil, 0, err
		}
		debts = append(debts, debt)
	}

	total, err := countDebts()
	if err != nil {
		return nil, 0, err
	}

	return debts, total, nil
}

func countDebts() (int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM debts").Scan(&total)
	return total, err
}
