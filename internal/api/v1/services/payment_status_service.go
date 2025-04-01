package services

import (
	"backend-go/internal/api/models"
	"backend-go/internal/api/v1/dto"
	"backend-go/internal/api/v1/repository"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
)

type PaymentStatusService struct {
	DB *repository.Database
}

func NewPaymentStatusService(db *repository.Database) *PaymentStatusService {
	return &PaymentStatusService{DB: db}
}

func (s *PaymentStatusService) ParsePaymentStatus(categoryReq dto.PaymentStatusRequest) (models.PaymentStatus, error) {
	return models.PaymentStatus{
		Name:        categoryReq.Name,
		Description: &categoryReq.Description,
	}, nil

}

func (s *PaymentStatusService) CreatePaymentStatus(PaymentStatus models.PaymentStatus) (models.PaymentStatus, error) {
	return s.DB.InsertPaymentStatus(PaymentStatus)
}

func (s *PaymentStatusService) UpdatePaymentStatus(PaymentStatus models.PaymentStatus) (models.PaymentStatus, error) {
	return s.DB.UpdatePaymentStatus(PaymentStatus)
}

func (s *PaymentStatusService) ListPaymentStatus(pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, int, error) {
	invoices, err := s.DB.ListPaymentStatus(pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountPaymentStatus(pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (s *PaymentStatusService) GetPaymentStatusByID(id uuid.UUID) (*models.PaymentStatus, error) {
	return s.DB.GetPaymentStatusByID(id)
}

func (s *PaymentStatusService) DeletePaymentStatusByID(id uuid.UUID) error {
	return s.DB.DeletePaymentStatusByID(id)
}
