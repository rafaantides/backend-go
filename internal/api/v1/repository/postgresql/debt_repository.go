package postgresql

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/ent"

	// "backend-go/pkg/pagination"
	"context"
	"fmt"

	// "strings"
	"time"

	"github.com/google/uuid"
	// "github.com/lib/pq"
)

func (d *PostgreSQL) GetDebtByID(ctx context.Context, id uuid.UUID) (models.Debt, error) {
	row, err := d.Client.Debt.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return models.Debt{}, errs.ErrNotFound
		}
		return models.Debt{}, err
	}

	data, err := newDebtResponse(row)
	if err != nil {
		return models.Debt{}, fmt.Errorf("failed to insert debt: %w", err)
	}
	return data, nil
}

// func (d *PostgreSQL) DeleteDebtByID(id uuid.UUID) error {
// 	query := `DELETE FROM debts WHERE id = $1`
// 	result, err := d.DB.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return errs.ErrNotFound
// 	}
// 	return nil
// }

// func (d *PostgreSQL) InsertDebt(debt models.Debt) (models.Debt, error) {
// 	query := `INSERT INTO debts (invoice_id, title, category_id, amount, purchase_date, due_date)
// 			  VALUES ($1, $2, $3, $4, $5, $6)
// 			  RETURNING *`

// 	row := d.DB.QueryRow(query, debt.InvoiceID, debt.Title, debt.CategoryID,
// 		debt.Amount, debt.PurchaseDate, debt.DueDate)
// 	data, err := newDebtResponse(row)
// 	if err != nil {
// 		return models.Debt{}, fmt.Errorf("failed to insert debt: %w", err)
// 	}
// 	return data, nil*sql.Rows
// }

// func (d *PostgreSQL) UpdateDebt(debt models.Debt) (models.Debt, error) {
// 	query := `
// 		UPDATE debts
// 		SET title = $1, amount = $2, purchase_date = $3, category_id = $4 , status_id = $5
// 		WHERE id = $6
// 		RETURNING *
// 	`
// 	row := d.DB.QueryRow(query, debt.Title, debt.Amount, debt.PurchaseDate, debt.CategoryID, debt.StatusID, debt.ID)
// 	data, err := newDebtResponse(row)
// 	if err != nil {
// 		return models.Debt{}, fmt.Errorf("failed to update debt: %w", err)
// 	}
// 	return data, nil
// }

// func (d *PostgreSQL) ListDebts(flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtsResponse, error) {
// 	query := `
//         SELECT
//             d.id,
//             d.title,
//             d.amount,
//             d.purchase_date,
//             d.due_date,
//             d.category_id,
//             d.status_id,
//             d.created_at,
//             d.updated_at,
//             c.name AS category,
//             i.title AS invoice_title,
//             s.name AS status
//         FROM debts d
//         LEFT JOIN categories c ON d.category_id = c.id
//         LEFT JOIN payment_status s ON d.status_id = s.id
//         LEFT JOIN invoices i ON d.invoice_id = i.id
//     `
// 	filterQuery, args := buildDebtFilters(flt, pgn)
// 	query += filterQuery
// 	query += fmt.Sprintf(" ORDER BY d.%s DESC", pgn.OrderBy)
// 	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
// 	args = append(args, pgn.PageSize, pgn.Offset())

// 	rows, err := d.DB.Query(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return newDebtsResponse(rows)
// }

// func (d *PostgreSQL) CountDebts(flt dto.DebtFilters, pgn *pagination.Pagination) (int, error) {
// 	query := "SELECT COUNT(*) FROM debts d LEFT JOIN categories c ON d.category_id = c.id LEFT JOIN payment_status s ON d.status_id = s.id LEFT JOIN invoices i ON d.invoice_id = i.id"
// 	filterQuery, args := buildDebtFilters(flt, pgn)
// 	query += filterQuery

// 	var total int
// 	err := d.DB.QueryRow(query, args...).Scan(&total)
// 	return total, err
// }

func newDebtResponse(debt *ent.Debt) (models.Debt, error) {
	if debt == nil {
		return models.Debt{}, fmt.Errorf("debt is nil")
	}

	return models.Debt{
		ID:           debt.ID,
		Title:        debt.Title,
		InvoiceID:    &debt.Edges.Invoice.ID,
		CategoryID:   &debt.Edges.Category.ID,
		Amount:       debt.Amount,
		PurchaseDate: debt.PurchaseDate,
		DueDate:      debt.DueDate,
		StatusID:     &debt.Edges.Status.ID,
		CreatedAt:    debt.CreatedAt,
		UpdatedAt:    debt.UpdatedAt,
	}, nil
}

// TODO: mudar para o utils
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func newDebtsResponse(debts []*ent.Debt) ([]dto.DebtsResponse, error) {
	if debts == nil {
		return nil, fmt.Errorf("debts slice is nil")
	}

	response := make([]dto.DebtsResponse, 0, len(debts))
	for _, debt := range debts {
		response = append(response, dto.DebtsResponse{
			ID:     debt.ID,
			Title:  debt.Title,
			Amount: debt.Amount,
			// TODO: colocar invoiceID
			// TODO: mudar para um utils
			PurchaseDate: debt.PurchaseDate.Format("2006-01-02"),
			DueDate:      strPtr(debt.DueDate.Format("2006-01-02")),
			CategoryID:   &debt.Edges.Category.ID,
			StatusID:     &debt.Edges.Status.ID,
			// TODO: mudar para um utils
			CreatedAt:    debt.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    debt.UpdatedAt.Format(time.RFC3339),
			Category:     &debt.Edges.Category.Name,
			InvoiceTitle: &debt.Edges.Invoice.Title,
			Status:       &debt.Edges.Status.Name,
		})
	}

	return response, nil
}

// func buildDebtFilters(flt dto.DebtFilters, pgn *pagination.Pagination) (string, []any) {
// 	var conditions []string
// 	var args []any
// 	argIndex := 1

// 	if pgn.Search != "" {
// 		conditions = append(conditions, fmt.Sprintf(
// 			"(d.title ILIKE $%d OR c.name ILIKE $%d OR i.title ILIKE $%d OR s.name ILIKE $%d)",
// 			argIndex, argIndex+1, argIndex+2, argIndex+3,
// 		))
// 		args = append(args, "%"+pgn.Search+"%", "%"+pgn.Search+"%", "%"+pgn.Search+"%", "%"+pgn.Search+"%")
// 		argIndex += 4
// 	}
// 	if flt.CategoryID != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.category_id = ANY($%d)", argIndex))
// 		args = append(args, pq.Array(*flt.CategoryID))
// 		argIndex++
// 	}
// 	if flt.StatusID != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.status_id = ANY($%d)", argIndex))
// 		args = append(args, pq.Array(*flt.StatusID))
// 		argIndex++
// 	}
// 	if flt.InvoiceID != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.invoice_id = ANY($%d)", argIndex))
// 		args = append(args, pq.Array(*flt.InvoiceID))
// 		argIndex++
// 	}
// 	if flt.MinAmount != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.amount >= $%d", argIndex))
// 		args = append(args, *flt.MinAmount)
// 		argIndex++
// 	}
// 	if flt.MaxAmount != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.amount <= $%d", argIndex))
// 		args = append(args, *flt.MaxAmount)
// 		argIndex++
// 	}
// 	if flt.StartDate != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.purchase_date >= $%d", argIndex))
// 		args = append(args, *flt.StartDate)
// 		argIndex++
// 	}
// 	if flt.EndDate != nil {
// 		conditions = append(conditions, fmt.Sprintf("d.purchase_date <= $%d", argIndex))
// 		args = append(args, *flt.EndDate)
// 		argIndex++
// 	}

// 	filterQuery := ""
// 	if len(conditions) > 0 {
// 		filterQuery = " WHERE " + strings.Join(conditions, " AND ")
// 	}

// 	return filterQuery, args
// }
