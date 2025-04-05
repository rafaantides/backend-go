package postgresql

import (
	"backend-go/internal/api/errs"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/ent"
	"backend-go/pkg/ent/category"
	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

func (d *PostgreSQL) GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	row, err := d.Client.Category.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return newCategoryResponse(row)
}

func (d *PostgreSQL) GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error) {
	if name == nil {
		return nil, nil
	}

	data, err := d.Client.Category.Query().Where(category.NameEQ(*name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	id := data.ID
	return &id, nil
}

func (d *PostgreSQL) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	err := d.Client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errs.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) InsertCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error) {
	created, err := d.Client.Category.
		Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)

	if err != nil {
		return nil, errs.FailedToSave("categories", err)
	}

	return newCategoryResponse(created)
}

func (d *PostgreSQL) UpdateCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error) {
	updated, err := d.Client.Category.
		UpdateOneID(input.ID).
		SetName(input.Name).
		SetNillableDescription(input.Description).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.FailedToSave("categories", err)
	}

	return newCategoryResponse(updated)
}

func (d *PostgreSQL) ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, error) {
	query := d.Client.Category.Query()

	query = applyCategoryFilters(query, pgn)
	query = query.Order(ent.Desc(pgn.OrderBy))
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newCategoryResponseList(data)
}

func (d *PostgreSQL) CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error) {
	query := d.Client.Category.Query()
	query = applyCategoryFilters(query, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func mapCategoryToResponse(row *ent.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
	}
}

func newCategoryResponse(row *ent.Category) (*dto.CategoryResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapCategoryToResponse(row)
	return &response, nil
}

func newCategoryResponseList(rows []*ent.Category) ([]dto.CategoryResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.CategoryResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapCategoryToResponse(row))
	}
	return response, nil
}

func applyCategoryFilters(query *ent.CategoryQuery, pgn *pagination.Pagination) *ent.CategoryQuery {
	if pgn.Search != "" {
		query = query.Where(
			category.Or(
				category.NameContainsFold(pgn.Search),
				category.DescriptionContainsFold(pgn.Search),
			),
		)
	}
	return query
}
