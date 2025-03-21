package repository

import (
	"api-go/internal/errs"
	"api-go/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var model models.Category

	query := `SELECT * FROM categories WHERE id = $1`
	row := DB.QueryRow(query, id)

	err := row.Scan(&model.ID, &model.Name, &model.Description, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNoRows
		}
		return nil, err
	}

	return &model, nil
}

func GetCategoryIDByName(categoryName *string) (*uuid.UUID, error) {
	if categoryName == nil {
		return nil, nil
	}
	query := "SELECT id FROM categories WHERE name = $1"

	var categoryID uuid.UUID
	err := DB.QueryRow(query, categoryName).Scan(&categoryID)
	if err != nil {
		// TODO: rever se é melhor retornar nil ou um erro
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &categoryID, nil
}
