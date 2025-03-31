package services

import (
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
)

func ParsePaymentStatus(categoryReq dto.PaymentStatusRequest) (models.PaymentStatus, error) {
	return models.PaymentStatus{
		Name:        categoryReq.Name,
		Description: &categoryReq.Description,
	}, nil

}

func CreatePaymentStatus(PaymentStatus models.PaymentStatus) (models.PaymentStatus, error) {
	return repository.InsertPaymentStatus(PaymentStatus)
}

func UpdatePaymentStatus(PaymentStatus models.PaymentStatus) (models.PaymentStatus, error) {
	return repository.UpdatePaymentStatus(PaymentStatus)
}

func ListPaymentStatus(pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, int, error) {
	invoices, err := repository.ListPaymentStatus(pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := repository.CountPaymentStatus(pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func GetPaymentStatusByID(id uuid.UUID) (*models.PaymentStatus, error) {
	return repository.GetPaymentStatusByID(id)
}

func DeletePaymentStatusByID(id uuid.UUID) error {
	return repository.DeletePaymentStatusByID(id)
}
