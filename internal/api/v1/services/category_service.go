package services

import (
	"backend-go/internal/api/v1/dto"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

type CategoryService struct {
	DB repository.Database
}

func NewCategoryService(db repository.Database) *CategoryService {
	return &CategoryService{DB: db}
}

func (s *CategoryService) ParseCategory(req dto.CategoryRequest) (models.Category, error) {
	return models.Category{
		Name:        req.Name,
		Description: &req.Description,
	}, nil

}

func (s *CategoryService) CreateCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error) {
	return s.DB.InsertCategory(ctx, input)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, input models.Category) (*dto.CategoryResponse, error) {
	return s.DB.UpdateCategory(ctx, input)
}

func (s *CategoryService) ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, int, error) {
	data, err := s.DB.ListCategories(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountCategories(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	return s.DB.GetCategoryByID(ctx, id)
}

func (s *CategoryService) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeleteCategoryByID(ctx, id)
}
