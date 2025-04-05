package postgresql

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/ent"
	"backend-go/pkg/ent/invoice"
	"backend-go/pkg/ent/paymentstatus"
	"backend-go/pkg/pagination"
	"backend-go/pkg/utils"
	"context"

	"github.com/google/uuid"
)

func (d *PostgreSQL) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceResponse, error) {
	row, err := d.Client.Invoice.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return newInvoiceResponse(row)
}

func (d *PostgreSQL) DeleteInvoiceByID(ctx context.Context, id uuid.UUID) error {
	err := d.Client.Invoice.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errs.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) InsertInvoice(ctx context.Context, input models.Invoice) (*dto.InvoiceResponse, error) {
	created, err := d.Client.Invoice.
		Create().
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetIssueDate(input.IssueDate).
		SetDueDate(input.DueDate).
		Save(ctx)

	if err != nil {
		return nil, errs.FailedToSave("invoices", err)
	}
	return newInvoiceResponse(created)
}

func (d *PostgreSQL) UpdateInvoice(ctx context.Context, input models.Invoice) (*dto.InvoiceResponse, error) {
	updated, err := d.Client.Invoice.
		UpdateOneID(input.ID).
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetIssueDate(input.IssueDate).
		SetDueDate(input.DueDate).
		SetStatusID(input.StatusID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.FailedToSave("invoices", err)
	}
	return newInvoiceResponse(updated)
}

func (d *PostgreSQL) ListInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error) {
	query := d.Client.Invoice.Query()

	query = applyInvoiceFilters(query, flt, pgn)
	query = query.Order(ent.Desc(pgn.OrderBy))
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newInvoiceResponseList(data)
}

func (d *PostgreSQL) CountInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error) {
	query := d.Client.Invoice.Query()
	query = applyInvoiceFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func mapInvoiceToResponse(row *ent.Invoice) dto.InvoiceResponse {
	var statusID *uuid.UUID
	var statusName *string

	if row.Edges.Status != nil {
		statusID = &row.Edges.Status.ID
		statusName = &row.Edges.Status.Name
	}

	return dto.InvoiceResponse{
		ID:        row.ID,
		Title:     row.Title,
		Amount:    row.Amount,
		IssueDate: *utils.ToFormatDatePointer(row.IssueDate),
		DueDate:   utils.ToFormatDatePointer(row.DueDate),
		CreatedAt: *utils.ToFormatDateTimePointer(row.CreatedAt),
		UpdatedAt: *utils.ToFormatDateTimePointer(row.UpdatedAt),
		StatusID:  statusID,
		Status:    statusName,
	}
}

func newInvoiceResponse(row *ent.Invoice) (*dto.InvoiceResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapInvoiceToResponse(row)
	return &response, nil
}

func newInvoiceResponseList(rows []*ent.Invoice) ([]dto.InvoiceResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.InvoiceResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapInvoiceToResponse(row))
	}
	return response, nil
}

func applyInvoiceFilters(query *ent.InvoiceQuery, flt dto.InvoiceFilters, pgn *pagination.Pagination) *ent.InvoiceQuery {
	if pgn.Search != "" {
		query = query.Where(
			invoice.Or(
				invoice.TitleContainsFold(pgn.Search),
				invoice.HasStatusWith(
					paymentstatus.NameContainsFold(pgn.Search),
				),
			),
		)
	}
	if flt.StatusID != nil {
		statusIds := utils.ToUUIDSlice(*flt.StatusID)
		if len(statusIds) > 0 {
			query = query.Where(
				invoice.HasStatusWith(paymentstatus.IDIn(statusIds...)),
			)
		}
	}
	if flt.MinAmount != nil {
		query = query.Where(
			invoice.AmountGTE(*flt.MinAmount),
		)
	}
	if flt.MaxAmount != nil {
		query = query.Where(
			invoice.AmountLTE(*flt.MaxAmount),
		)
	}
	if t := utils.ToTimePointer(flt.StartDate); t != nil {
		query = query.Where(invoice.IssueDateGTE(*t))
	}

	if t := utils.ToTimePointer(flt.EndDate); t != nil {
		query = query.Where(invoice.IssueDateLTE(*t))
	}

	return query
}
