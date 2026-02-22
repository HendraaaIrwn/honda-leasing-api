package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type KecamatanService interface {
	CRUDService[models.Kecamatan]
	GetAllKecamatan(ctx context.Context) ([]models.Kecamatan, error)
	GetKecamatanByID(ctx context.Context, id int64) (*models.Kecamatan, error)
	CreateKecamatan(ctx context.Context, kecamatan *models.Kecamatan) error
	UpdateKecamatan(ctx context.Context, kecamatan *models.Kecamatan) error
	DeleteKecamatan(ctx context.Context, id int64) error
	ListByKabID(ctx context.Context, kabID int64) ([]models.Kecamatan, error)
}

type kecamatanService struct {
	*baseService[models.Kecamatan]
	repo repository.KecamatanRepository
}

func NewKecamatanService(repo repository.KecamatanRepository) KecamatanService {
	return &kecamatanService{
		baseService: newBaseService[models.Kecamatan](repo),
		repo:        repo,
	}
}

func (s *kecamatanService) GetAllKecamatan(ctx context.Context) ([]models.Kecamatan, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *kecamatanService) GetKecamatanByID(ctx context.Context, id int64) (*models.Kecamatan, error) {
	if err := validatePositiveID("kecamatan", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *kecamatanService) CreateKecamatan(ctx context.Context, kecamatan *models.Kecamatan) error {
	if kecamatan == nil {
		return errors.New("kecamatan payload cannot be empty")
	}
	if err := validateName("kecamatan", kecamatan.KecName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("kabupaten", kecamatan.KabID); err != nil {
		return err
	}
	return s.repo.Create(ctx, kecamatan)
}

func (s *kecamatanService) UpdateKecamatan(ctx context.Context, kecamatan *models.Kecamatan) error {
	if kecamatan == nil {
		return errors.New("kecamatan payload cannot be empty")
	}
	if err := validatePositiveID("kecamatan", kecamatan.KecID); err != nil {
		return err
	}
	if err := validateName("kecamatan", kecamatan.KecName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("kabupaten", kecamatan.KabID); err != nil {
		return err
	}
	return s.repo.Update(ctx, kecamatan.KecID, map[string]interface{}{
		"kec_name": strings.TrimSpace(kecamatan.KecName),
		"kab_id":   kecamatan.KabID,
	})
}

func (s *kecamatanService) DeleteKecamatan(ctx context.Context, id int64) error {
	if err := validatePositiveID("kecamatan", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *kecamatanService) ListByKabID(ctx context.Context, kabID int64) ([]models.Kecamatan, error) {
	if err := validatePositiveID("kabupaten", kabID); err != nil {
		return nil, err
	}
	return s.repo.ListByKabID(ctx, kabID)
}
