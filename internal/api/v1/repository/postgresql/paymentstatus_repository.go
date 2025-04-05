package postgresql

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/ent"
	"backend-go/pkg/ent/paymentstatus"
	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

func (d *PostgreSQL) GetPaymentStatusByID(ctx context.Context, id uuid.UUID) (*models.PaymentStatus, error) {
	row, err := d.Client.PaymentStatus.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return newPaymentStatusResponse(row)
}

func (d *PostgreSQL) GetPaymentStatusIDByName(ctx context.Context, name *string) (*uuid.UUID, error) {
	if name == nil {
		return nil, nil
	}

	data, err := d.Client.PaymentStatus.Query().Where(paymentstatus.NameEQ(*name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	id := data.ID
	return &id, nil
}

func (d *PostgreSQL) DeletePaymentStatusByID(ctx context.Context, id uuid.UUID) error {
	err := d.Client.PaymentStatus.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errs.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) InsertPaymentStatus(ctx context.Context, input models.PaymentStatus) (*models.PaymentStatus, error) {
	created, err := d.Client.PaymentStatus.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)

	if err != nil {
		return nil, errs.FailedToSave("payment_status", err)
	}

	return newPaymentStatusResponse(created)
}

func (d *PostgreSQL) UpdatePaymentStatus(ctx context.Context, input models.PaymentStatus) (*models.PaymentStatus, error) {
	updated, err := d.Client.PaymentStatus.
		UpdateOneID(input.ID).
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.FailedToSave("payment_status", err)
	}

	return newPaymentStatusResponse(updated)
}

func (d *PostgreSQL) ListPaymentStatus(ctx context.Context, pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, error) {
	query := d.Client.PaymentStatus.Query()

	query = applyPaymentStatusFilters(query, pgn)
	query = query.Order(ent.Desc(pgn.OrderBy))
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newPaymentsStatusResponse(data)
}

func (d *PostgreSQL) CountPaymentStatus(ctx context.Context, pgn *pagination.Pagination) (int, error) {
	query := d.Client.PaymentStatus.Query()
	query = applyPaymentStatusFilters(query, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func newPaymentStatusResponse(row *ent.PaymentStatus) (*models.PaymentStatus, error) {
	return &models.PaymentStatus{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}, nil
}

func newPaymentsStatusResponse(rows []*ent.PaymentStatus) ([]dto.PaymentStatusResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.PaymentStatusResponse, len(rows))
	for i, row := range rows {
		response[i] = dto.PaymentStatusResponse{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
		}
	}
	return response, nil
}

func applyPaymentStatusFilters(query *ent.PaymentStatusQuery, pgn *pagination.Pagination) *ent.PaymentStatusQuery {
	if pgn.Search != "" {
		query = query.Where(
			paymentstatus.Or(
				paymentstatus.NameContainsFold(pgn.Search),
				paymentstatus.DescriptionContainsFold(pgn.Search),
			),
		)
	}
	return query
}
