package repository

import (
	"context"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gorm"
)

type MotorTypeRepository interface {
	CRUDRepository[models.MotorType]
	GetByName(ctx context.Context, name string) (*models.MotorType, error)
}

type MotorRepository interface {
	CRUDRepository[models.Motor]
	GetByNomorPolisi(ctx context.Context, nomorPolisi string) (*models.Motor, error)
	GetByNomorRangka(ctx context.Context, nomorRangka string) (*models.Motor, error)
}

type MotorAssetRepository interface {
	CRUDRepository[models.MotorAsset]
	ListByMotorID(ctx context.Context, motorID int64) ([]models.MotorAsset, error)
}

type CustomerRepository interface {
	CRUDRepository[models.Customer]
	GetByNIK(ctx context.Context, nik string) (*models.Customer, error)
	GetByEmail(ctx context.Context, email string) (*models.Customer, error)
}

type motorTypeRepository struct {
	*baseRepository[models.MotorType]
}

type motorRepository struct {
	*baseRepository[models.Motor]
}

type motorAssetRepository struct {
	*baseRepository[models.MotorAsset]
}

type customerRepository struct {
	*baseRepository[models.Customer]
}

func NewMotorTypeRepository(db *gorm.DB) MotorTypeRepository {
	return &motorTypeRepository{baseRepository: newBaseRepository[models.MotorType](db)}
}

func NewMotorRepository(db *gorm.DB) MotorRepository {
	return &motorRepository{baseRepository: newBaseRepository[models.Motor](db)}
}

func NewMotorAssetRepository(db *gorm.DB) MotorAssetRepository {
	return &motorAssetRepository{baseRepository: newBaseRepository[models.MotorAsset](db)}
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{baseRepository: newBaseRepository[models.Customer](db)}
}

func (r *motorTypeRepository) GetByName(ctx context.Context, name string) (*models.MotorType, error) {
	value, err := validateLookupValue(name)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "moty_name = ?", value)
}

func (r *motorRepository) GetByNomorPolisi(ctx context.Context, nomorPolisi string) (*models.Motor, error) {
	value, err := validateLookupValue(nomorPolisi)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "nomor_polisi = ?", value)
}

func (r *motorRepository) GetByNomorRangka(ctx context.Context, nomorRangka string) (*models.Motor, error) {
	value, err := validateLookupValue(nomorRangka)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "nomor_rangka = ?", value)
}

func (r *motorAssetRepository) ListByMotorID(ctx context.Context, motorID int64) ([]models.MotorAsset, error) {
	if motorID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.MotorAsset
	if err := r.db.WithContext(ctx).Where("moas_motor_id = ?", motorID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *customerRepository) GetByNIK(ctx context.Context, nik string) (*models.Customer, error) {
	value, err := validateLookupValue(nik)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "nik = ?", value)
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	value, err := validateLookupValue(email)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "email = ?", value)
}
