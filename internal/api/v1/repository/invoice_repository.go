package repository

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/pkg/pagination"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func GetInvoiceByID(id uuid.UUID) (*models.Invoice, error) {
	var data models.Invoice

	query := `SELECT * FROM invoices WHERE id = $1`
	row := DB.QueryRow(query, id)

	err := row.Scan(&data.ID, &data.Title, &data.Amount, &data.IssueDate, &data.DueDate, &data.StatusID, &data.CreatedAt, &data.UpdatedAt)
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

func InsertInvoice(invoice models.Invoice) (models.Invoice, error) {
	query := `INSERT INTO debts (title, amount, issue_date, due_date)
			  VALUES ($1, $2, $3, $4)
			  RETURNING *`

	var newData models.Invoice
	err := DB.QueryRow(query, invoice.Title, invoice.Amount, invoice.IssueDate, invoice.DueDate).
		Scan(&newData.ID, &newData.Title, &newData.Amount, &newData.IssueDate, &newData.DueDate, &newData.StatusID, &newData.CreatedAt, &newData.UpdatedAt)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to insert invoice: %w", err)
	}
	return newData, nil
}

func UpdateInvoice(invoice models.Invoice) (models.Invoice, error) {
	query := `
		UPDATE invoices
		SET title = $1, amount = $2, issue_date = $3, due_date = $4, status_id = $5
		WHERE id = $6
		RETURNING *
	`
	var updatedData models.Invoice
	err := DB.QueryRow(query, invoice.Title, invoice.Amount, invoice.IssueDate, invoice.DueDate, invoice.StatusID, invoice.ID).
		Scan(&updatedData.ID, &updatedData.Title, &updatedData.Amount, &updatedData.IssueDate, &updatedData.DueDate, &updatedData.StatusID, &updatedData.CreatedAt, &updatedData.UpdatedAt)

	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to update debt: %w", err)
	}

	return updatedData, nil
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

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += fmt.Sprintf(" ORDER BY i.%s DESC", pgn.OrderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pgn.PageSize, pgn.Offset())

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return newInvoiceResponse(rows)
}

func CountInvoices() (int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&total)
	return total, err
}

func newInvoiceResponse(rows *sql.Rows) ([]dto.InvoiceResponse, error) {
	defer rows.Close()
	invoices := make([]dto.InvoiceResponse, 0)
	for rows.Next() {
		var invoice dto.InvoiceResponse

		err := rows.Scan(
			&invoice.ID, &invoice.Title, &invoice.Amount, &invoice.IssueDate, &invoice.DueDate,
			&invoice.StatusID, &invoice.CreatedAt, &invoice.UpdatedAt, &invoice.Status,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
