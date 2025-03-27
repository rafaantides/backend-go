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

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	row := DB.QueryRow(`SELECT * FROM categories WHERE id = $1`, id)
	data, err := newCategoryResponse(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &data, nil
}

func GetCategoryIDByName(categoryName *string) (*uuid.UUID, error) {
	if categoryName == nil {
		return nil, nil
	}
	var categoryID uuid.UUID
	err := DB.QueryRow(`SELECT id FROM categories WHERE name = $1`, categoryName).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &categoryID, nil
}

func DeleteCategoryByID(id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
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

func InsertCategory(input models.Category) (models.Category, error) {
	query := `INSERT INTO categories (name, description)
			  VALUES ($1, $2)
			  RETURNING id, name, description`

	row := DB.QueryRow(query, input.Name, input.Description)
	data, err := newCategoryResponse(row)
	if err != nil {
		return models.Category{}, errs.FailedToSave("categories", err)
	}

	return data, nil
}

func UpdateCategory(input models.Category) (models.Category, error) {
	query := `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING *
	`

	row := DB.QueryRow(query, input.Name, input.Description)
	data, err := newCategoryResponse(row)
	if err != nil {
		return models.Category{}, errs.FailedToSave("categories", err)
	}
	return data, nil
}

func ListCategories(pgn *pagination.Pagination) ([]dto.CategoryResponse, error) {
	query := `
        SELECT
			id,
            name,
			description
		FROM categories
    `

	filterQuery, args := buildCategoryFilters(pgn)
	query += filterQuery

	argIndex := len(args) + 1
	query += fmt.Sprintf(" ORDER BY %s DESC", pgn.OrderBy)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pgn.PageSize, pgn.Offset())

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return newCategoriesResponse(rows)
}

func CountCategories(pgn *pagination.Pagination) (int, error) {
	query := "SELECT COUNT(*) FROM categories"
	filterQuery, args := buildCategoryFilters(pgn)
	query += filterQuery

	var total int
	err := DB.QueryRow(query, args...).Scan(&total)
	return total, err
}

func newCategoryResponse(row *sql.Row) (models.Category, error) {

	var data models.Category
	if err := row.Scan(&data.ID, &data.Name, &data.Description); err != nil {
		return models.Category{}, err
	}

	return data, nil
}

func newCategoriesResponse(rows *sql.Rows) ([]dto.CategoryResponse, error) {
	defer rows.Close()
	categories := make([]dto.CategoryResponse, 0)
	for rows.Next() {
		var category dto.CategoryResponse

		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func buildCategoryFilters(pgn *pagination.Pagination) (string, []any) {
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
