package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type KelurahanService interface {
	CRUDService[models.Kelurahan]
	GetAllKelurahan(ctx context.Context) ([]models.Kelurahan, error)
	GetKelurahanByID(ctx context.Context, id int64) (*models.Kelurahan, error)
	CreateKelurahan(ctx context.Context, kelurahan *models.Kelurahan) error
	UpdateKelurahan(ctx context.Context, kelurahan *models.Kelurahan) error
	DeleteKelurahan(ctx context.Context, id int64) error
	ListByKecID(ctx context.Context, kecID int64) ([]models.Kelurahan, error)
}

type kelurahanService struct {
	*baseService[models.Kelurahan]
	repo repository.KelurahanRepository
}

func NewKelurahanService(repo repository.KelurahanRepository) KelurahanService {
	return &kelurahanService{
		baseService: newBaseService[models.Kelurahan](repo),
		repo:        repo,
	}
}

func (s *kelurahanService) GetAllKelurahan(ctx context.Context) ([]models.Kelurahan, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *kelurahanService) GetKelurahanByID(ctx context.Context, id int64) (*models.Kelurahan, error) {
	if err := validatePositiveID("kelurahan", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *kelurahanService) CreateKelurahan(ctx context.Context, kelurahan *models.Kelurahan) error {
	if kelurahan == nil {
		return errors.New("kelurahan payload cannot be empty")
	}
	if err := validateName("kelurahan", kelurahan.KelName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("kecamatan", kelurahan.KecID); err != nil {
		return err
	}
	return s.repo.Create(ctx, kelurahan)
}

func (s *kelurahanService) UpdateKelurahan(ctx context.Context, kelurahan *models.Kelurahan) error {
	if kelurahan == nil {
		return errors.New("kelurahan payload cannot be empty")
	}
	if err := validatePositiveID("kelurahan", kelurahan.KelID); err != nil {
		return err
	}
	if err := validateName("kelurahan", kelurahan.KelName, maxRegionNameLength); err != nil {
		return err
	}
	if err := validatePositiveID("kecamatan", kelurahan.KecID); err != nil {
		return err
	}
	return s.repo.Update(ctx, kelurahan.KelID, map[string]interface{}{
		"kel_name": strings.TrimSpace(kelurahan.KelName),
		"kec_id":   kelurahan.KecID,
	})
}

func (s *kelurahanService) DeleteKelurahan(ctx context.Context, id int64) error {
	if err := validatePositiveID("kelurahan", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *kelurahanService) ListByKecID(ctx context.Context, kecID int64) ([]models.Kelurahan, error) {
	if err := validatePositiveID("kecamatan", kecID); err != nil {
		return nil, err
	}
	return s.repo.ListByKecID(ctx, kecID)
}
