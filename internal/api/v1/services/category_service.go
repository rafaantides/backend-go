package services

import (
	"backend-go/internal/api/v1/dto"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
)

type CategoryService struct {
	DB repository.Database
}

func NewCategoryService(db repository.Database) *CategoryService {
	return &CategoryService{DB: db}
}

func (s *CategoryService) ParseCategory(categoryReq dto.CategoryRequest) (models.Category, error) {
	return models.Category{
		Name:        categoryReq.Name,
		Description: &categoryReq.Description,
	}, nil

}

func (s *CategoryService) CreateCategory(input models.Category) (models.Category, error) {
	return s.DB.InsertCategory(input)
}

func (s *CategoryService) UpdateCategory(input models.Category) (models.Category, error) {
	return s.DB.UpdateCategory(input)
}

func (s *CategoryService) ListCategories(pgn *pagination.Pagination) ([]dto.CategoriesResponse, int, error) {
	data, err := s.DB.ListCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountCategories(pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return s.DB.GetCategoryByID(id)
}

func (s *CategoryService) DeleteCategoryByID(id uuid.UUID) error {
	return s.DB.DeleteCategoryByID(id)
}
