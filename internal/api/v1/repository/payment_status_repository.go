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
)

func (d *Database) GetPaymentStatusByID(id uuid.UUID) (*models.PaymentStatus, error) {
	row := d.DB.QueryRow(`SELECT * FROM payment_status WHERE id = $1`, id)
	data, err := newPaymentStatusResponse(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &data, nil
}

func (d *Database) GetPaymentStatusIDByName(name *string) (*uuid.UUID, error) {
	if name == nil {
		return nil, nil
	}
	var id uuid.UUID
	err := d.DB.QueryRow(`SELECT id FROM payment_status WHERE name = $1`, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &id, nil
}

func (d *Database) DeletePaymentStatusByID(id uuid.UUID) error {
	query := `DELETE FROM payment_status WHERE id = $1`
	result, err := d.DB.Exec(query, id)
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

func (d *Database) InsertPaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error) {
	query := `INSERT INTO payment_status (name, description)
			  VALUES ($1, $2)
			  RETURNING id, name, description`

	row := d.DB.QueryRow(query, input.Name, input.Description)
	data, err := newPaymentStatusResponse(row)
	if err != nil {
		return models.PaymentStatus{}, errs.FailedToSave("payment_status", err)
	}

	return data, nil
}

func (d *Database) UpdatePaymentStatus(input models.PaymentStatus) (models.PaymentStatus, error) {
	query := `
		UPDATE payment_status
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING *
	`

	row := d.DB.QueryRow(query, input.Name, input.Description)
	data, err := newPaymentStatusResponse(row)
	if err != nil {
		return models.PaymentStatus{}, errs.FailedToSave("payment_status", err)
	}
	return data, nil
}

func (d *Database) ListPaymentStatus(pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, error) {
	query := `
        SELECT
			id,
            name,
			description
		FROM payment_status
    `

	filterQuery, args := buildPaymentStatusFilters(pgn)
	query += filterQuery

	argIndex := len(args) + 1
	query += fmt.Sprintf(" ORDER BY %s DESC", pgn.OrderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pgn.PageSize, pgn.Offset())

	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return newPaymentsStatusResponse(rows)
}

func (d *Database) CountPaymentStatus(pgn *pagination.Pagination) (int, error) {
	query := "SELECT COUNT(*) FROM payment_status"
	filterQuery, args := buildPaymentStatusFilters(pgn)
	query += filterQuery

	var total int
	err := d.DB.QueryRow(query, args...).Scan(&total)
	return total, err
}

func newPaymentStatusResponse(row *sql.Row) (models.PaymentStatus, error) {
	var data models.PaymentStatus
	if err := row.Scan(&data.ID, &data.Name, &data.Description); err != nil {
		return models.PaymentStatus{}, err
	}
	return data, nil
}

func newPaymentsStatusResponse(rows *sql.Rows) ([]dto.PaymentStatusResponse, error) {
	defer rows.Close()
	response := make([]dto.PaymentStatusResponse, 0)
	for rows.Next() {
		var data dto.PaymentStatusResponse
		if err := rows.Scan(&data.ID, &data.Name, &data.Description); err != nil {
			return make([]dto.PaymentStatusResponse, 0), err
		}
		response = append(response, data)
	}
	return response, nil
}

func buildPaymentStatusFilters(pgn *pagination.Pagination) (string, []any) {
	var conditions []string
	var args []any
	argIndex := 1

	if pgn.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(name ILIKE $%d OR description ILIKE $%d)",
			argIndex, argIndex+1,
		))
		args = append(args, "%"+pgn.Search+"%", "%"+pgn.Search+"%")
		argIndex += 2
	}

	filterQuery := ""
	if len(conditions) > 0 {
		filterQuery = " WHERE " + strings.Join(conditions, " AND ")
	}

	return filterQuery, args
}
