package repository

import (
	"context"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gorm"
)

type ProvinceRepository interface {
	CRUDRepository[models.Province]
	GetByName(ctx context.Context, name string) (*models.Province, error)
}

type KabupatenRepository interface {
	CRUDRepository[models.Kabupaten]
	ListByProvID(ctx context.Context, provID int64) ([]models.Kabupaten, error)
}

type KecamatanRepository interface {
	CRUDRepository[models.Kecamatan]
	ListByKabID(ctx context.Context, kabID int64) ([]models.Kecamatan, error)
}

type KelurahanRepository interface {
	CRUDRepository[models.Kelurahan]
	ListByKecID(ctx context.Context, kecID int64) ([]models.Kelurahan, error)
}

type LocationRepository interface {
	CRUDRepository[models.Location]
	ListByKelID(ctx context.Context, kelID int64) ([]models.Location, error)
}

type TemplateTaskRepository interface {
	CRUDRepository[models.TemplateTask]
	ListByRoleID(ctx context.Context, roleID int64) ([]models.TemplateTask, error)
}

type TemplateTaskAttributeRepository interface {
	CRUDRepository[models.TemplateTaskAttribute]
	ListByTaskID(ctx context.Context, taskID int64) ([]models.TemplateTaskAttribute, error)
}

type provinceRepository struct {
	*baseRepository[models.Province]
}

type kabupatenRepository struct {
	*baseRepository[models.Kabupaten]
}

type kecamatanRepository struct {
	*baseRepository[models.Kecamatan]
}

type kelurahanRepository struct {
	*baseRepository[models.Kelurahan]
}

type locationRepository struct {
	*baseRepository[models.Location]
}

type templateTaskRepository struct {
	*baseRepository[models.TemplateTask]
}

type templateTaskAttributeRepository struct {
	*baseRepository[models.TemplateTaskAttribute]
}

func NewProvinceRepository(db *gorm.DB) ProvinceRepository {
	return &provinceRepository{baseRepository: newBaseRepository[models.Province](db)}
}

func NewKabupatenRepository(db *gorm.DB) KabupatenRepository {
	return &kabupatenRepository{baseRepository: newBaseRepository[models.Kabupaten](db)}
}

func NewKecamatanRepository(db *gorm.DB) KecamatanRepository {
	return &kecamatanRepository{baseRepository: newBaseRepository[models.Kecamatan](db)}
}

func NewKelurahanRepository(db *gorm.DB) KelurahanRepository {
	return &kelurahanRepository{baseRepository: newBaseRepository[models.Kelurahan](db)}
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{baseRepository: newBaseRepository[models.Location](db)}
}

func NewTemplateTaskRepository(db *gorm.DB) TemplateTaskRepository {
	return &templateTaskRepository{baseRepository: newBaseRepository[models.TemplateTask](db)}
}

func NewTemplateTaskAttributeRepository(db *gorm.DB) TemplateTaskAttributeRepository {
	return &templateTaskAttributeRepository{baseRepository: newBaseRepository[models.TemplateTaskAttribute](db)}
}

func (r *provinceRepository) GetByName(ctx context.Context, name string) (*models.Province, error) {
	value, err := validateLookupValue(name)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "prov_name = ?", value)
}

func (r *kabupatenRepository) ListByProvID(ctx context.Context, provID int64) ([]models.Kabupaten, error) {
	if provID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Kabupaten
	if err := r.db.WithContext(ctx).Where("prov_id = ?", provID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *kecamatanRepository) ListByKabID(ctx context.Context, kabID int64) ([]models.Kecamatan, error) {
	if kabID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Kecamatan
	if err := r.db.WithContext(ctx).Where("kab_id = ?", kabID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *kelurahanRepository) ListByKecID(ctx context.Context, kecID int64) ([]models.Kelurahan, error) {
	if kecID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Kelurahan
	if err := r.db.WithContext(ctx).Where("kec_id = ?", kecID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *locationRepository) ListByKelID(ctx context.Context, kelID int64) ([]models.Location, error) {
	if kelID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.Location
	if err := r.db.WithContext(ctx).Where("kel_id = ?", kelID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *templateTaskRepository) ListByRoleID(ctx context.Context, roleID int64) ([]models.TemplateTask, error) {
	if roleID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.TemplateTask
	if err := r.db.WithContext(ctx).Where("teta_role_id = ?", roleID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *templateTaskAttributeRepository) ListByTaskID(ctx context.Context, taskID int64) ([]models.TemplateTaskAttribute, error) {
	if taskID < 1 {
		return nil, errs.ErrInvalidInput
	}

	var items []models.TemplateTaskAttribute
	if err := r.db.WithContext(ctx).Where("tetat_teta_id = ?", taskID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
