package services

import (
	"context"
	"errors"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type LocationService interface {
	CRUDService[models.Location]
	GetAllLocations(ctx context.Context) ([]models.Location, error)
	GetLocationByID(ctx context.Context, id int64) (*models.Location, error)
	CreateLocation(ctx context.Context, location *models.Location) error
	UpdateLocation(ctx context.Context, location *models.Location) error
	DeleteLocation(ctx context.Context, id int64) error
	ListByKelID(ctx context.Context, kelID int64) ([]models.Location, error)
}

type locationService struct {
	*baseService[models.Location]
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{
		baseService: newBaseService[models.Location](repo),
		repo:        repo,
	}
}

func (s *locationService) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	return fetchAllRecords(ctx, s.repo)
}

func (s *locationService) GetLocationByID(ctx context.Context, id int64) (*models.Location, error) {
	if err := validatePositiveID("location", id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *locationService) CreateLocation(ctx context.Context, location *models.Location) error {
	if location == nil {
		return errors.New("location payload cannot be empty")
	}
	if err := validatePositiveID("kelurahan", location.KelID); err != nil {
		return err
	}
	if err := validatePostalCode(location.PostalCode); err != nil {
		return err
	}
	return s.repo.Create(ctx, location)
}

func (s *locationService) UpdateLocation(ctx context.Context, location *models.Location) error {
	if location == nil {
		return errors.New("location payload cannot be empty")
	}
	if err := validatePositiveID("location", location.LocationID); err != nil {
		return err
	}
	if err := validatePositiveID("kelurahan", location.KelID); err != nil {
		return err
	}
	if err := validatePostalCode(location.PostalCode); err != nil {
		return err
	}
	return s.repo.Update(ctx, location.LocationID, map[string]interface{}{
		"street_address": strings.TrimSpace(location.StreetAddress),
		"postal_code":    strings.TrimSpace(location.PostalCode),
		"longitude":      strings.TrimSpace(location.Longitude),
		"latitude":       strings.TrimSpace(location.Latitude),
		"kel_id":         location.KelID,
	})
}

func (s *locationService) DeleteLocation(ctx context.Context, id int64) error {
	if err := validatePositiveID("location", id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *locationService) ListByKelID(ctx context.Context, kelID int64) ([]models.Location, error) {
	if err := validatePositiveID("kelurahan", kelID); err != nil {
		return nil, err
	}
	return s.repo.ListByKelID(ctx, kelID)
}
