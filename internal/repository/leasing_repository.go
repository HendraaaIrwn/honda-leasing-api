package repository

import (
	"context"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gorm"
)

type LeasingProductRepository interface {
	CRUDRepository[models.LeasingProduct]
	GetByKodeProduk(ctx context.Context, kodeProduk string) (*models.LeasingProduct, error)
}

type LeasingContractRepository interface {
	CRUDRepository[models.LeasingContract]
	GetByContractNumber(ctx context.Context, contractNumber string) (*models.LeasingContract, error)
	ListByCustomerID(ctx context.Context, customerID int64) ([]models.LeasingContract, error)
}

type LeasingTaskRepository interface {
	CRUDRepository[models.LeasingTask]
	ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingTask, error)
}

type LeasingTaskAttributeRepository interface {
	CRUDRepository[models.LeasingTaskAttribute]
	ListByTaskID(ctx context.Context, taskID int64) ([]models.LeasingTaskAttribute, error)
}

type LeasingContractDocumentRepository interface {
	CRUDRepository[models.LeasingContractDocument]
	ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingContractDocument, error)
}

type leasingProductRepository struct {
	*baseRepository[models.LeasingProduct]
}

type leasingContractRepository struct {
	*baseRepository[models.LeasingContract]
}

type leasingTaskRepository struct {
	*baseRepository[models.LeasingTask]
}

type leasingTaskAttributeRepository struct {
	*baseRepository[models.LeasingTaskAttribute]
}

type leasingContractDocumentRepository struct {
	*baseRepository[models.LeasingContractDocument]
}

func NewLeasingProductRepository(db *gorm.DB) LeasingProductRepository {
	return &leasingProductRepository{baseRepository: newBaseRepository[models.LeasingProduct](db)}
}

func NewLeasingContractRepository(db *gorm.DB) LeasingContractRepository {
	return &leasingContractRepository{baseRepository: newBaseRepository[models.LeasingContract](db)}
}

func NewLeasingTaskRepository(db *gorm.DB) LeasingTaskRepository {
	return &leasingTaskRepository{baseRepository: newBaseRepository[models.LeasingTask](db)}
}

func NewLeasingTaskAttributeRepository(db *gorm.DB) LeasingTaskAttributeRepository {
	return &leasingTaskAttributeRepository{baseRepository: newBaseRepository[models.LeasingTaskAttribute](db)}
}

func NewLeasingContractDocumentRepository(db *gorm.DB) LeasingContractDocumentRepository {
	return &leasingContractDocumentRepository{baseRepository: newBaseRepository[models.LeasingContractDocument](db)}
}

func (r *leasingProductRepository) GetByKodeProduk(ctx context.Context, kodeProduk string) (*models.LeasingProduct, error) {
	value, err := validateLookupValue(kodeProduk)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "kode_produk = ?", value)
}

func (r *leasingContractRepository) GetByContractNumber(ctx context.Context, contractNumber string) (*models.LeasingContract, error) {
	value, err := validateLookupValue(contractNumber)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "contract_number = ?", value)
}

func (r *leasingContractRepository) ListByCustomerID(ctx context.Context, customerID int64) ([]models.LeasingContract, error) {
	if customerID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.LeasingContract
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *leasingTaskRepository) ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingTask, error) {
	if contractID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.LeasingTask
	if err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *leasingTaskAttributeRepository) ListByTaskID(ctx context.Context, taskID int64) ([]models.LeasingTaskAttribute, error) {
	if taskID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.LeasingTaskAttribute
	if err := r.db.WithContext(ctx).Where("tasa_leta_id = ?", taskID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *leasingContractDocumentRepository) ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingContractDocument, error) {
	if contractID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.LeasingContractDocument
	if err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
