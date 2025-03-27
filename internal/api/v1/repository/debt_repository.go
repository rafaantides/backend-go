package repository

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/pkg/pagination"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func newDebtResponse(row *sql.Row) (models.Debt, error) {
	var data models.Debt
	if err := row.Scan(
		&data.ID, &data.Title, &data.CategoryID, &data.Amount, &data.PurchaseDate,
		&data.DueDate, &data.StatusID, &data.CreatedAt, &data.UpdatedAt,
	); err != nil {
		return models.Debt{}, err
	}

	return data, nil
}

func newDebtsResponse(rows *sql.Rows) ([]dto.DebtResponse, error) {
	defer rows.Close()
	response := make([]dto.DebtResponse, 0)
	for rows.Next() {
		var data dto.DebtResponse
		if err := rows.Scan(
			&data.ID, &data.Title, &data.Amount, &data.PurchaseDate, &data.DueDate,
			&data.CategoryID, &data.StatusID, &data.CreatedAt, &data.UpdatedAt,
			&data.Category, &data.InvoiceTitle, &data.Status,
		); err != nil {
			return make([]dto.DebtResponse, 0), err
		}
		response = append(response, data)
	}
	return response, nil
}

func GetDebtByID(id uuid.UUID) (*models.Debt, error) {
	row := DB.QueryRow(`SELECT * FROM debts WHERE id = $1`, id)
	data, err := newDebtResponse(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
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
		return errs.ErrNotFound
	}
	return nil
}

func InsertDebt(debt models.Debt) (models.Debt, error) {
	query := `INSERT INTO debts (invoice_id, title, category_id, amount, purchase_date, due_date)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING *`

	row := DB.QueryRow(query, debt.InvoiceID, debt.Title, debt.CategoryID,
		debt.Amount, debt.PurchaseDate, debt.DueDate)
	data, err := newDebtResponse(row)
	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to insert debt: %w", err)
	}
	return data, nil
}

func UpdateDebt(debt models.Debt) (models.Debt, error) {
	query := `
		UPDATE debts 
		SET title = $1, amount = $2, purchase_date = $3, category_id = $4 , status_id = $5
		WHERE id = $6
		RETURNING *
	`
	row := DB.QueryRow(query, debt.Title, debt.Amount, debt.PurchaseDate, debt.CategoryID, debt.StatusID, debt.ID)
	data, err := newDebtResponse(row)
	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to update debt: %w", err)
	}
	return data, nil
}

func ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, error) {
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
	filterQuery, args := buildDebtFilters(flt, pgn)
	query += filterQuery
	query += fmt.Sprintf(" ORDER BY d.%s DESC", pgn.OrderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, pgn.PageSize, pgn.Offset())

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return newDebtsResponse(rows)
}

func CountDebts(flt dto.DebtFilters, pgn *pagination.Pagination) (int, error) {
	query := "SELECT COUNT(*) FROM debts d LEFT JOIN categories c ON d.category_id = c.id LEFT JOIN payment_status s ON d.status_id = s.id LEFT JOIN invoices i ON d.invoice_id = i.id"
	filterQuery, args := buildDebtFilters(flt, pgn)
	query += filterQuery

	var total int
	err := DB.QueryRow(query, args...).Scan(&total)
	return total, err
}

func buildDebtFilters(flt dto.DebtFilters, pgn *pagination.Pagination) (string, []any) {
	var conditions []string
	var args []any
	argIndex := 1

	if pgn.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(d.title ILIKE $%d OR c.name ILIKE $%d OR i.title ILIKE $%d OR s.name ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3,
		))
		args = append(args, "%"+pgn.Search+"%", "%"+pgn.Search+"%", "%"+pgn.Search+"%", "%"+pgn.Search+"%")
		argIndex += 4
	}
	if flt.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("d.category_id = ANY($%d)", argIndex))
		args = append(args, pq.Array(*flt.CategoryID))
		argIndex++
	}
	if flt.StatusID != nil {
		conditions = append(conditions, fmt.Sprintf("d.status_id = ANY($%d)", argIndex))
		args = append(args, pq.Array(*flt.StatusID))
		argIndex++
	}
	if flt.InvoiceID != nil {
		conditions = append(conditions, fmt.Sprintf("d.invoice_id = ANY($%d)", argIndex))
		args = append(args, pq.Array(*flt.InvoiceID))
		argIndex++
	}
	if flt.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount >= $%d", argIndex))
		args = append(args, *flt.MinAmount)
		argIndex++
	}
	if flt.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("d.amount <= $%d", argIndex))
		args = append(args, *flt.MaxAmount)
		argIndex++
	}
	if flt.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date >= $%d", argIndex))
		args = append(args, *flt.StartDate)
		argIndex++
	}
	if flt.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("d.purchase_date <= $%d", argIndex))
		args = append(args, *flt.EndDate)
		argIndex++
	}

	filterQuery := ""
	if len(conditions) > 0 {
		filterQuery = " WHERE " + strings.Join(conditions, " AND ")
	}

	return filterQuery, args
}
