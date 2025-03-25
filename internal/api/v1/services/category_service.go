package services

import (
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/repository"

	"github.com/google/uuid"
)

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return repository.GetCategoryByID(id)
}
