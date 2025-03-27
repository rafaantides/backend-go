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

func CreateCategory(Category models.Category) (models.Category, error) {
	return repository.InsertCategory(Category)
}

func UpdateCategory(Category models.Category) (models.Category, error) {
	return repository.UpdateCategory(Category)
}

func ListCategories(pgn *pagination.Pagination) ([]dto.CategoryResponse, int, error) {
	invoices, err := repository.ListCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := repository.CountCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return repository.GetCategoryByID(id)
}

func DeleteCategoryByID(id uuid.UUID) error {
	return repository.DeleteCategoryByID(id)
}
