package postgresql

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/ent"
	"backend-go/pkg/ent/category"
	"backend-go/pkg/ent/debt"
	"backend-go/pkg/ent/invoice"
	"backend-go/pkg/ent/paymentstatus"
	"backend-go/pkg/utils"

	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

func (d *PostgreSQL) GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error) {
	row, err := d.Client.Debt.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return newDebtResponse(row)
}

func (d *PostgreSQL) DeleteDebtByID(ctx context.Context, id uuid.UUID) error {
	err := d.Client.Debt.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errs.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) InsertDebt(ctx context.Context, input models.Debt) (*dto.DebtResponse, error) {
	created, err := d.Client.Debt.
		Create().
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetDueDate(input.DueDate).
		SetPurchaseDate(input.PurchaseDate).
		SetStatusID(*input.StatusID).
		SetInvoiceID(*input.InvoiceID).
		SetCategoryID(*input.CategoryID).
		Save(ctx)

	if err != nil {
		return nil, errs.FailedToSave("debts", err)
	}
	return newDebtResponse(created)
}

func (d *PostgreSQL) UpdateDebt(ctx context.Context, input models.Debt) (*dto.DebtResponse, error) {
	updated, err := d.Client.Debt.
		UpdateOneID(input.ID).
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetDueDate(input.DueDate).
		SetPurchaseDate(input.PurchaseDate).
		SetStatusID(*input.StatusID).
		SetInvoiceID(*input.InvoiceID).
		SetCategoryID(*input.CategoryID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.FailedToSave("debts", err)
	}
	return newDebtResponse(updated)
}

func (d *PostgreSQL) ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, error) {
	query := d.Client.Debt.Query().
		WithStatus().
		WithCategory().
		WithInvoice()

	query = applyDebtFilters(query, flt, pgn)
	query = query.Order(ent.Desc(pgn.OrderBy))
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newDebtResponseList(data)
}

func (d *PostgreSQL) CountDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) (int, error) {
	query := d.Client.Debt.Query()
	query = applyDebtFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func mapDebtToResponse(row *ent.Debt) dto.DebtResponse {
	var categoryID *uuid.UUID
	var categoryName *string
	var invoiceID *uuid.UUID
	var invoiceTitle *string
	var statusID *uuid.UUID
	var statusName *string

	if row.Edges.Category != nil {
		categoryID = &row.Edges.Category.ID
		categoryName = &row.Edges.Category.Name
	}

	if row.Edges.Invoice != nil {
		invoiceID = &row.Edges.Invoice.ID
		invoiceTitle = &row.Edges.Invoice.Title
	}

	if row.Edges.Status != nil {
		statusID = &row.Edges.Status.ID
		statusName = &row.Edges.Status.Name
	}

	return dto.DebtResponse{
		ID:           row.ID,
		Title:        row.Title,
		Amount:       row.Amount,
		PurchaseDate: *utils.ToFormatDateTimePointer(row.PurchaseDate),
		DueDate:      utils.ToFormatDatePointer(row.DueDate),
		CategoryID:   categoryID,
		Category:     categoryName,
		StatusID:     statusID,
		Status:       statusName,
		CreatedAt:    *utils.ToFormatDateTimePointer(row.CreatedAt),
		UpdatedAt:    *utils.ToFormatDateTimePointer(row.UpdatedAt),
		InvoiceID:    invoiceID,
		InvoiceTitle: invoiceTitle,
	}
}

func newDebtResponse(row *ent.Debt) (*dto.DebtResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapDebtToResponse(row)
	return &response, nil
}

func newDebtResponseList(rows []*ent.Debt) ([]dto.DebtResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.DebtResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapDebtToResponse(row))
	}
	return response, nil
}

func applyDebtFilters(query *ent.DebtQuery, flt dto.DebtFilters, pgn *pagination.Pagination) *ent.DebtQuery {
	if pgn.Search != "" {
		query = query.Where(
			debt.Or(
				debt.TitleContainsFold(pgn.Search),
				debt.HasStatusWith(
					paymentstatus.NameContainsFold(pgn.Search),
				),
				debt.HasCategoryWith(
					category.NameContainsFold(pgn.Search),
				),
				debt.HasInvoiceWith(
					invoice.TitleContains(pgn.Search),
				),
			),
		)
	}

	if flt.StatusID != nil {
		statusIds := utils.ToUUIDSlice(*flt.StatusID)
		if len(statusIds) > 0 {
			query = query.Where(
				debt.HasStatusWith(paymentstatus.IDIn(statusIds...)),
			)
		}
	}
	if flt.CategoryID != nil {
		categoryIds := utils.ToUUIDSlice(*flt.CategoryID)
		if len(categoryIds) > 0 {
			query = query.Where(
				debt.HasCategoryWith(category.IDIn(categoryIds...)),
			)
		}
	}
	if flt.InvoiceID != nil {
		invoiceIds := utils.ToUUIDSlice(*flt.InvoiceID)
		if len(invoiceIds) > 0 {
			query = query.Where(
				debt.HasInvoiceWith(invoice.IDIn(invoiceIds...)),
			)
		}
	}
	if flt.MinAmount != nil {
		query = query.Where(
			debt.AmountGTE(*flt.MinAmount),
		)
	}
	if flt.MaxAmount != nil {
		query = query.Where(
			debt.AmountLTE(*flt.MaxAmount),
		)
	}
	if t := utils.ToTimePointer(flt.StartDate); t != nil {
		query = query.Where(debt.PurchaseDateGTE(*t))
	}

	if t := utils.ToTimePointer(flt.EndDate); t != nil {
		query = query.Where(debt.PurchaseDateLTE(*t))
	}

	return query
}
