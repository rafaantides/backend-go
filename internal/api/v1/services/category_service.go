package services

import (
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
)

func ParseCategory(categoryReq dto.CategoryRequest) (models.Category, error) {
	return models.Category{
		Name:        categoryReq.Name,
		Description: &categoryReq.Description,
	}, nil

}

func CreateCategory(input models.Category) (models.Category, error) {
	return repository.InsertCategory(input)
}

func UpdateCategory(input models.Category) (models.Category, error) {
	return repository.UpdateCategory(input)
}

func ListCategories(pgn *pagination.Pagination) ([]dto.CategoriesResponse, int, error) {
	data, err := repository.ListCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := repository.CountCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return repository.GetCategoryByID(id)
}

func DeleteCategoryByID(id uuid.UUID) error {
	return repository.DeleteCategoryByID(id)
}
