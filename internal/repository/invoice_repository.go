package repository

import (
	"api-go/internal/dto"
	"api-go/internal/errs"
	"api-go/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func GetInvoiceByID(id uuid.UUID) (*models.Invoice, error) {
	var data models.Invoice

	query := `SELECT * FROM invoices WHERE id = $1`
	row := DB.QueryRow(query, id)

	err := row.Scan(&data.ID, &data.Title, &data.Amount, &data.IssueDate, &data.DueDate, &data.StatusID, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNoRows
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
		return errs.ErrNoRows
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

func ListInvoices(filters dto.InvoiceFilters, page int, pageSize int, orderBy string) ([]dto.InvoiceResponse, int, error) {
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

	if filters.Title != nil {
		conditions = append(conditions, fmt.Sprintf("i.title ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Title+"%")
		argIndex++
	}
	if filters.StatusID != nil {
		conditions = append(conditions, fmt.Sprintf("i.status_id = $%d", argIndex))
		args = append(args, *filters.StatusID)
		argIndex++
	}
	if filters.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("i.amount >= $%d", argIndex))
		args = append(args, *filters.MinAmount)
		argIndex++
	}
	if filters.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("i.amount <= $%d", argIndex))
		args = append(args, *filters.MaxAmount)
		argIndex++
	}
	if filters.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("i.issue_date >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}
	if filters.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("i.issue_date <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += fmt.Sprintf(" ORDER BY i.%s DESC", orderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []dto.InvoiceResponse
	invoices = make([]dto.InvoiceResponse, 0)
	for rows.Next() {
		var invoice dto.InvoiceResponse
		err := rows.Scan(
			&invoice.ID, &invoice.Title, &invoice.Amount, &invoice.IssueDate, &invoice.DueDate,
			&invoice.StatusID, &invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, invoice)
	}

	total, err := countInvoices()
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func countInvoices() (int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&total)
	return total, err
}
