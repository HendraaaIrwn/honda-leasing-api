package services

import (
	"context"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type LeasingProductService interface {
	CRUDService[models.LeasingProduct]
	GetByKodeProduk(ctx context.Context, kodeProduk string) (*models.LeasingProduct, error)
}

type LeasingContractService interface {
	CRUDService[models.LeasingContract]
	GetByContractNumber(ctx context.Context, contractNumber string) (*models.LeasingContract, error)
	ListByCustomerID(ctx context.Context, customerID int64) ([]models.LeasingContract, error)
}

type LeasingTaskService interface {
	CRUDService[models.LeasingTask]
	ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingTask, error)
}

type LeasingTaskAttributeService interface {
	CRUDService[models.LeasingTaskAttribute]
	ListByTaskID(ctx context.Context, taskID int64) ([]models.LeasingTaskAttribute, error)
}

type LeasingContractDocumentService interface {
	CRUDService[models.LeasingContractDocument]
	ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingContractDocument, error)
}

type leasingProductService struct {
	*baseService[models.LeasingProduct]
	repo repository.LeasingProductRepository
}

type leasingContractService struct {
	*baseService[models.LeasingContract]
	repo repository.LeasingContractRepository
}

type leasingTaskService struct {
	*baseService[models.LeasingTask]
	repo repository.LeasingTaskRepository
}

type leasingTaskAttributeService struct {
	*baseService[models.LeasingTaskAttribute]
	repo repository.LeasingTaskAttributeRepository
}

type leasingContractDocumentService struct {
	*baseService[models.LeasingContractDocument]
	repo repository.LeasingContractDocumentRepository
}

func NewLeasingProductService(repo repository.LeasingProductRepository) LeasingProductService {
	return &leasingProductService{
		baseService: newBaseService[models.LeasingProduct](repo),
		repo:        repo,
	}
}

func NewLeasingContractService(repo repository.LeasingContractRepository) LeasingContractService {
	return &leasingContractService{
		baseService: newBaseService[models.LeasingContract](repo),
		repo:        repo,
	}
}

func NewLeasingTaskService(repo repository.LeasingTaskRepository) LeasingTaskService {
	return &leasingTaskService{
		baseService: newBaseService[models.LeasingTask](repo),
		repo:        repo,
	}
}

func NewLeasingTaskAttributeService(repo repository.LeasingTaskAttributeRepository) LeasingTaskAttributeService {
	return &leasingTaskAttributeService{
		baseService: newBaseService[models.LeasingTaskAttribute](repo),
		repo:        repo,
	}
}

func NewLeasingContractDocumentService(repo repository.LeasingContractDocumentRepository) LeasingContractDocumentService {
	return &leasingContractDocumentService{
		baseService: newBaseService[models.LeasingContractDocument](repo),
		repo:        repo,
	}
}

func (s *leasingProductService) GetByKodeProduk(ctx context.Context, kodeProduk string) (*models.LeasingProduct, error) {
	return s.repo.GetByKodeProduk(ctx, kodeProduk)
}

func (s *leasingContractService) GetByContractNumber(ctx context.Context, contractNumber string) (*models.LeasingContract, error) {
	return s.repo.GetByContractNumber(ctx, contractNumber)
}

func (s *leasingContractService) ListByCustomerID(ctx context.Context, customerID int64) ([]models.LeasingContract, error) {
	return s.repo.ListByCustomerID(ctx, customerID)
}

func (s *leasingTaskService) ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingTask, error) {
	return s.repo.ListByContractID(ctx, contractID)
}

func (s *leasingTaskAttributeService) ListByTaskID(ctx context.Context, taskID int64) ([]models.LeasingTaskAttribute, error) {
	return s.repo.ListByTaskID(ctx, taskID)
}

func (s *leasingContractDocumentService) ListByContractID(ctx context.Context, contractID int64) ([]models.LeasingContractDocument, error) {
	return s.repo.ListByContractID(ctx, contractID)
}
