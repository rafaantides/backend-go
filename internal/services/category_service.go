package services

import (
	"api-go/internal/models"
	"api-go/internal/repository"

	"github.com/google/uuid"
)

func GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return repository.GetCategoryByID(id)
}
