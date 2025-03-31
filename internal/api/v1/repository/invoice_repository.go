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

func newInvoiceResponse(row *sql.Row) (models.Invoice, error) {
	var data models.Invoice
	if err := row.Scan(
		&data.ID, &data.Title, &data.Amount, &data.IssueDate, &data.DueDate,
		&data.StatusID, &data.CreatedAt, &data.UpdatedAt,
	); err != nil {
		return models.Invoice{}, err
	}
	return data, nil
}

func newInvoicesResponse(rows *sql.Rows) ([]dto.InvoiceResponse, error) {
	defer rows.Close()
	response := make([]dto.InvoiceResponse, 0)
	for rows.Next() {
		var data dto.InvoiceResponse
		if err := rows.Scan(
			&data.ID, &data.Title, &data.Amount, &data.IssueDate, &data.DueDate,
			&data.StatusID, &data.CreatedAt, &data.UpdatedAt, &data.Status,
		); err != nil {
			return make([]dto.InvoiceResponse, 0), err
		}
		response = append(response, data)
	}
	return response, nil
}

func GetInvoiceByID(id uuid.UUID) (*models.Invoice, error) {
	row := DB.QueryRow(`SELECT * FROM invoices WHERE id = $1`, id)
	data, err := newInvoiceResponse(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &data, nil
}

func DeleteInvoiceByID(id uuid.UUID) error {
	query := `DELETE FROM invoices WHERE id = $1`
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

func InsertInvoice(input models.Invoice) (models.Invoice, error) {
	query := `INSERT INTO debts (title, amount, issue_date, due_date)
			  VALUES ($1, $2, $3, $4)
			  RETURNING *`

	row := DB.QueryRow(query, input.Title, input.Amount, input.IssueDate, input.DueDate)
	data, err := newInvoiceResponse(row)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to insert invoice: %w", err)
	}
	return data, nil
}

func UpdateInvoice(input models.Invoice) (models.Invoice, error) {
	query := `
		UPDATE invoices
		SET title = $1, amount = $2, issue_date = $3, due_date = $4, status_id = $5
		WHERE id = $6
		RETURNING *
	`
	row := DB.QueryRow(query, input.Title, input.Amount, input.IssueDate, input.DueDate, input.StatusID, input.ID)
	data, err := newInvoiceResponse(row)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to update debt: %w", err)
	}
	return data, nil
}

func ListInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error) {
	query := `
        SELECT
            i.id,
			i.title,
			i.amount,
			i.issue_date,
			i.due_date,
			i.status_id,
            i.created_at,
			i.updated_at,
			s.name AS status
		FROM invoices i
		LEFT JOIN payment_status s ON i.status_id = s.id
    `

	filterQuery, args := buildInvoiceFilters(flt, pgn)
	query += filterQuery

	argIndex := len(args) + 1
	query += fmt.Sprintf(" ORDER BY i.%s DESC", pgn.OrderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pgn.PageSize, pgn.Offset())

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return newInvoicesResponse(rows)
}

func CountInvoices(flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error) {
	query := "SELECT COUNT(*) FROM invoices i LEFT JOIN payment_status s ON i.status_id = s.id"
	filterQuery, args := buildInvoiceFilters(flt, pgn)
	query += filterQuery

	var total int
	err := DB.QueryRow(query, args...).Scan(&total)
	return total, err
}

func buildInvoiceFilters(flt dto.InvoiceFilters, pgn *pagination.Pagination) (string, []any) {
	var conditions []string
	var args []any
	argIndex := 1

	if pgn.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(i.title ILIKE $%d OR s.name ILIKE $%d)",
			argIndex, argIndex+1,
		))
		args = append(args, "%"+pgn.Search+"%", "%"+pgn.Search+"%")
		argIndex += 2
	}
	if flt.StatusID != nil {
		conditions = append(conditions, fmt.Sprintf("i.status_id = ANY($%d)", argIndex))
		args = append(args, pq.Array(*flt.StatusID))
		argIndex++
	}
	if flt.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("i.amount >= $%d", argIndex))
		args = append(args, *flt.MinAmount)
		argIndex++
	}
	if flt.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("i.amount <= $%d", argIndex))
		args = append(args, *flt.MaxAmount)
		argIndex++
	}
	if flt.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("i.issue_date >= $%d", argIndex))
		args = append(args, *flt.StartDate)
		argIndex++
	}
	if flt.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("i.issue_date <= $%d", argIndex))
		args = append(args, *flt.EndDate)
		argIndex++
	}

	filterQuery := ""
	if len(conditions) > 0 {
		filterQuery = " WHERE " + strings.Join(conditions, " AND ")
	}

	return filterQuery, args
}
