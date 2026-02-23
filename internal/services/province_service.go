package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type ProvinceService interface {
	CRUDService[models.Province]
	GetAllProvinces(ctx context.Context) ([]models.Province, error)
	GetProvinceByID(ctx context.Context, id int64) (*models.Province, error)
	CreateProvince(ctx context.Context, province *models.Province) error
	UpdateProvince(ctx context.Context, province *models.Province) error
	DeleteProvince(ctx context.Context, id int64) error
	GetByName(ctx context.Context, name string) (*models.Province, error)
}

type provinceService struct {
	*baseService[models.Province]
	repo repository.ProvinceRepository
}

func NewProvinceService(repo repository.ProvinceRepository) ProvinceService {
	return &provinceService{
		baseService: newBaseService[models.Province](repo),
		repo:        repo,
	}
}

func (s *provinceService) GetAllProvinces(ctx context.Context) ([]models.Province, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *provinceService) GetProvinceByID(ctx context.Context, id int64) (*models.Province, error) {
	if err := validatePositiveID("province", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *provinceService) CreateProvince(ctx context.Context, province *models.Province) error {
	if province == nil {
		return errors.New("province payload cannot be empty")
	}
	if err := validateName("province", province.ProvName, maxRegionNameLength); err != nil {
		return err
	}
	return s.repo.Create(ctx, province)
}

func (s *provinceService) UpdateProvince(ctx context.Context, province *models.Province) error {
	if province == nil {
		return errors.New("province payload cannot be empty")
	}
	if err := validatePositiveID("province", province.ProvID); err != nil {
		return err
	}
	if err := validateName("province", province.ProvName, maxRegionNameLength); err != nil {
		return err
	}
	return s.repo.Update(ctx, province.ProvID, map[string]interface{}{
		"prov_name": strings.TrimSpace(province.ProvName),
	})
}

func (s *provinceService) DeleteProvince(ctx context.Context, id int64) error {
	if err := validatePositiveID("province", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *provinceService) GetByName(ctx context.Context, name string) (*models.Province, error) {
	if err := validateName("province", name, maxRegionNameLength); err != nil {
		return nil, err
	}
	return s.repo.GetByName(ctx, strings.TrimSpace(name))
}
