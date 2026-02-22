package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type KabupatenService interface {
	CRUDService[models.Kabupaten]
	GetAllKabupaten(ctx context.Context) ([]models.Kabupaten, error)
	GetKabupatenByID(ctx context.Context, id int64) (*models.Kabupaten, error)
	CreateKabupaten(ctx context.Context, kabupaten *models.Kabupaten) error
	UpdateKabupaten(ctx context.Context, kabupaten *models.Kabupaten) error
	DeleteKabupaten(ctx context.Context, id int64) error
	ListByProvID(ctx context.Context, provID int64) ([]models.Kabupaten, error)
}

type kabupatenService struct {
	*baseService[models.Kabupaten]
	repo repository.KabupatenRepository
}

func NewKabupatenService(repo repository.KabupatenRepository) KabupatenService {
	return &kabupatenService{
		baseService: newBaseService[models.Kabupaten](repo),
		repo:        repo,
	}
}

func (s *kabupatenService) GetAllKabupaten(ctx context.Context) ([]models.Kabupaten, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *kabupatenService) GetKabupatenByID(ctx context.Context, id int64) (*models.Kabupaten, error) {
	if err := validatePositiveID("kabupaten", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *kabupatenService) CreateKabupaten(ctx context.Context, kabupaten *models.Kabupaten) error {
	if kabupaten == nil {
		return errors.New("kabupaten payload cannot be empty")
	}
	if err := validateName("kabupaten", kabupaten.KabName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("province", kabupaten.ProvID); err != nil {
		return err
	}
	return s.repo.Create(ctx, kabupaten)
}

func (s *kabupatenService) UpdateKabupaten(ctx context.Context, kabupaten *models.Kabupaten) error {
	if kabupaten == nil {
		return errors.New("kabupaten payload cannot be empty")
	}
	if err := validatePositiveID("kabupaten", kabupaten.KabID); err != nil {
		return err
	}
	if err := validateName("kabupaten", kabupaten.KabName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("province", kabupaten.ProvID); err != nil {
		return err
	}
	return s.repo.Update(ctx, kabupaten.KabID, map[string]interface{}{
		"kab_name": strings.TrimSpace(kabupaten.KabName),
		"prov_id":  kabupaten.ProvID,
	})
}

func (s *kabupatenService) DeleteKabupaten(ctx context.Context, id int64) error {
	if err := validatePositiveID("kabupaten", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *kabupatenService) ListByProvID(ctx context.Context, provID int64) ([]models.Kabupaten, error) {
	if err := validatePositiveID("province", provID); err != nil {
		return nil, err
	}
	return s.repo.ListByProvID(ctx, provID)
}
