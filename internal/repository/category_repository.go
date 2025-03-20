package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

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
