package services

import (
	"context"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type MotorTypeService interface {
	CRUDService[models.MotorType]
	GetByName(ctx context.Context, name string) (*models.MotorType, error)
}

type MotorService interface {
	CRUDService[models.Motor]
	GetByNomorPolisi(ctx context.Context, nomorPolisi string) (*models.Motor, error)
	GetByNomorRangka(ctx context.Context, nomorRangka string) (*models.Motor, error)
}

type MotorAssetService interface {
	CRUDService[models.MotorAsset]
	ListByMotorID(ctx context.Context, motorID int64) ([]models.MotorAsset, error)
}

type CustomerService interface {
	CRUDService[models.Customer]
	GetByNIK(ctx context.Context, nik string) (*models.Customer, error)
	GetByEmail(ctx context.Context, email string) (*models.Customer, error)
}

type motorTypeService struct {
	*baseService[models.MotorType]
	repo repository.MotorTypeRepository
}

type motorService struct {
	*baseService[models.Motor]
	repo repository.MotorRepository
}

type motorAssetService struct {
	*baseService[models.MotorAsset]
	repo repository.MotorAssetRepository
}

type customerService struct {
	*baseService[models.Customer]
	repo repository.CustomerRepository
}

func NewMotorTypeService(repo repository.MotorTypeRepository) MotorTypeService {
	return &motorTypeService{
		baseService: newBaseService[models.MotorType](repo),
		repo:        repo,
	}
}

func NewMotorService(repo repository.MotorRepository) MotorService {
	return &motorService{
		baseService: newBaseService[models.Motor](repo),
		repo:        repo,
	}
}

func NewMotorAssetService(repo repository.MotorAssetRepository) MotorAssetService {
	return &motorAssetService{
		baseService: newBaseService[models.MotorAsset](repo),
		repo:        repo,
	}
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		baseService: newBaseService[models.Customer](repo),
		repo:        repo,
	}
}

func (s *motorTypeService) GetByName(ctx context.Context, name string) (*models.MotorType, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *motorService) GetByNomorPolisi(ctx context.Context, nomorPolisi string) (*models.Motor, error) {
	return s.repo.GetByNomorPolisi(ctx, nomorPolisi)
}

func (s *motorService) GetByNomorRangka(ctx context.Context, nomorRangka string) (*models.Motor, error) {
	return s.repo.GetByNomorRangka(ctx, nomorRangka)
}

func (s *motorAssetService) ListByMotorID(ctx context.Context, motorID int64) ([]models.MotorAsset, error) {
	return s.repo.ListByMotorID(ctx, motorID)
}

func (s *customerService) GetByNIK(ctx context.Context, nik string) (*models.Customer, error) {
	return s.repo.GetByNIK(ctx, nik)
}

func (s *customerService) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	return s.repo.GetByEmail(ctx, email)
}
