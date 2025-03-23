package services

import (
	"api-go/internal/config"
	"api-go/internal/dto"
	"api-go/internal/models"
	"api-go/internal/repository"
	"api-go/pkg/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ParseDebt(debtReq dto.DebtRequest) (models.Debt, error) {
	parsedDate, err := time.Parse("2006-01-02", debtReq.PurchaseDate)
	if err != nil {
		return models.Debt{}, fmt.Errorf("formato de data inválido, use YYYY-MM-DD: %w", err)
	}

	amount, err := strconv.ParseFloat(debtReq.Amount, 64)
	if err != nil {
		return models.Debt{}, fmt.Errorf("valor inválido: %w", err)
	}

	category := categorizeTransaction(debtReq.Title)
	categoryID, err := repository.GetCategoryIDByName(category)
	if err != nil {
		return models.Debt{}, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	invoiceID, err := utils.ToUUIDPointer(debtReq.InvoiceID)
	if err != nil {
		return models.Debt{}, fmt.Errorf("erro ao converter invoice_id: %w", err)
	}
	return models.Debt{
		InvoiceID:    invoiceID,
		Title:        debtReq.Title,
		Amount:       amount,
		PurchaseDate: parsedDate,
		CategoryID:   categoryID,
	}, nil

}

func categorizeTransaction(name string) *string {
	if category, exists := config.CategoryMap[name]; exists {
		return &category
	}
	return nil
}

func CreateDebt(debt models.Debt) (models.Debt, error) {
	return repository.InsertDebt(debt)
}

func UpdateDebt(debt models.Debt) (models.Debt, error) {
	return repository.UpdateDebt(debt)
}

func GetAllDebts(filters dto.DebtFilters, page int, pageSize int) ([]dto.DebtResponse, int, error) {
	return repository.GetAllDebts(filters, page, pageSize)
}

func GetDebtByID(id uuid.UUID) (*models.Debt, error) {
	return repository.GetDebtByID(id)
}

func DeleteDebtByID(id uuid.UUID) error {
	return repository.DeleteDebtByID(id)
}
