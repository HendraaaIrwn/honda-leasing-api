package services

import (
	"context"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type PaymentScheduleService interface {
	CRUDService[models.PaymentSchedule]
	ListByContractID(ctx context.Context, contractID int64) ([]models.PaymentSchedule, error)
}

type PaymentService interface {
	CRUDService[models.Payment]
	GetByNomorBukti(ctx context.Context, nomorBukti string) (*models.Payment, error)
	ListByContractID(ctx context.Context, contractID int64) ([]models.Payment, error)
	ListByScheduleID(ctx context.Context, scheduleID int64) ([]models.Payment, error)
}

type paymentScheduleService struct {
	*baseService[models.PaymentSchedule]
	repo repository.PaymentScheduleRepository
}

type paymentService struct {
	*baseService[models.Payment]
	repo repository.PaymentRepository
}

func NewPaymentScheduleService(repo repository.PaymentScheduleRepository) PaymentScheduleService {
	return &paymentScheduleService{
		baseService: newBaseService[models.PaymentSchedule](repo),
		repo:        repo,
	}
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		baseService: newBaseService[models.Payment](repo),
		repo:        repo,
	}
}

func (s *paymentScheduleService) ListByContractID(ctx context.Context, contractID int64) ([]models.PaymentSchedule, error) {
	return s.repo.ListByContractID(ctx, contractID)
}

func (s *paymentService) GetByNomorBukti(ctx context.Context, nomorBukti string) (*models.Payment, error) {
	return s.repo.GetByNomorBukti(ctx, nomorBukti)
}

func (s *paymentService) ListByContractID(ctx context.Context, contractID int64) ([]models.Payment, error) {
	return s.repo.ListByContractID(ctx, contractID)
}

func (s *paymentService) ListByScheduleID(ctx context.Context, scheduleID int64) ([]models.Payment, error) {
	return s.repo.ListByScheduleID(ctx, scheduleID)
}
