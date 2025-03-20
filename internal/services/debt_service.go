package services

import (
	"api-go/config"
	"api-go/internal/models"
	"api-go/internal/repository"
	"api-go/pkg/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

func ParseDebt(debtReq models.DebtRequest) (models.Debt, error) {
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

	return models.Debt{
		InvoiceID:    utils.ToUUIDPointer(debtReq.InvoiceID),
		Title:        debtReq.Title,
		Amount:       amount,
		PurchaseDate: parsedDate,
		CategoryID:   categoryID,
	}, nil

}

func CreateDebt(debt models.Debt) (models.Debt, error) {
	log.Printf("Criando débito no serviço: %+v", debt)

	createdDebt, err := repository.InsertDebt(debt)

	if err != nil {
		return models.Debt{}, fmt.Errorf("erro ao inserir débito: %w", err)
	}

	return createdDebt, nil
}

func categorizeTransaction(name string) *string {
	if category, exists := config.CategoryMap[name]; exists {
		return &category
	}
	return nil
}

func GetAllDebts(filters models.DebtFilters) ([]models.DebtResponse, int, error) {
	return repository.GetAllDebts(filters)
}
