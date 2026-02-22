package repository

import (
	"context"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gorm"
)

type PaymentScheduleRepository interface {
	CRUDRepository[models.PaymentSchedule]
	ListByContractID(ctx context.Context, contractID int64) ([]models.PaymentSchedule, error)
}

type PaymentRepository interface {
	CRUDRepository[models.Payment]
	GetByNomorBukti(ctx context.Context, nomorBukti string) (*models.Payment, error)
	ListByContractID(ctx context.Context, contractID int64) ([]models.Payment, error)
	ListByScheduleID(ctx context.Context, scheduleID int64) ([]models.Payment, error)
}

type paymentScheduleRepository struct {
	*baseRepository[models.PaymentSchedule]
}

type paymentRepository struct {
	*baseRepository[models.Payment]
}

func NewPaymentScheduleRepository(db *gorm.DB) PaymentScheduleRepository {
	return &paymentScheduleRepository{baseRepository: newBaseRepository[models.PaymentSchedule](db)}
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{baseRepository: newBaseRepository[models.Payment](db)}
}

func (r *paymentScheduleRepository) ListByContractID(ctx context.Context, contractID int64) ([]models.PaymentSchedule, error) {
	if contractID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.PaymentSchedule
	if err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *paymentRepository) GetByNomorBukti(ctx context.Context, nomorBukti string) (*models.Payment, error) {
	value, err := validateLookupValue(nomorBukti)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "nomor_bukti = ?", value)
}

func (r *paymentRepository) ListByContractID(ctx context.Context, contractID int64) ([]models.Payment, error) {
	if contractID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Payment
	if err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *paymentRepository) ListByScheduleID(ctx context.Context, scheduleID int64) ([]models.Payment, error) {
	if scheduleID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Payment
	if err := r.db.WithContext(ctx).Where("schedule_id = ?", scheduleID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
