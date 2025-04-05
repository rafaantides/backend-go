package services

import (
	"backend-go/internal/api/v1/dto"
	repository "backend-go/internal/api/v1/repository/interfaces"
	"backend-go/internal/api/v1/repository/models"
	"backend-go/pkg/pagination"
	"context"

	"github.com/google/uuid"
)

type PaymentStatusService struct {
	DB repository.Database
}

func NewPaymentStatusService(db repository.Database) *PaymentStatusService {
	return &PaymentStatusService{DB: db}
}

func (s *PaymentStatusService) ParsePaymentStatus(req dto.PaymentStatusRequest) (models.PaymentStatus, error) {
	return models.PaymentStatus{
		Name:        req.Name,
		Description: &req.Description,
	}, nil

}

func (s *PaymentStatusService) CreatePaymentStatus(ctx context.Context, input models.PaymentStatus) (*dto.PaymentStatusResponse, error) {
	return s.DB.InsertPaymentStatus(ctx, input)
}

func (s *PaymentStatusService) UpdatePaymentStatus(ctx context.Context, input models.PaymentStatus) (*dto.PaymentStatusResponse, error) {
	return s.DB.UpdatePaymentStatus(ctx, input)
}

func (s *PaymentStatusService) ListPaymentStatus(ctx context.Context, pgn *pagination.Pagination) ([]dto.PaymentStatusResponse, int, error) {
	invoices, err := s.DB.ListPaymentStatus(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.DB.CountPaymentStatus(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (s *PaymentStatusService) GetPaymentStatusByID(ctx context.Context, id uuid.UUID) (*dto.PaymentStatusResponse, error) {
	return s.DB.GetPaymentStatusByID(ctx, id)
}

func (s *PaymentStatusService) DeletePaymentStatusByID(ctx context.Context, id uuid.UUID) error {
	return s.DB.DeletePaymentStatusByID(ctx, id)
}
