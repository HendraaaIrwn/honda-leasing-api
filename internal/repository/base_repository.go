package repository

import (
	"context"
	"strings"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"gorm.io/gorm"
)

// CRUDRepository defines common operations for all entities.
type CRUDRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int64, preloads ...string) (*T, error)
	FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
	List(ctx context.Context, opts ListOptions, preloads ...string) ([]T, int64, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func newBaseRepository[T any](db *gorm.DB) *baseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) withPreloads(query *gorm.DB, preloads []string) *gorm.DB {
	for _, preload := range preloads {
		preload = strings.TrimSpace(preload)
		if preload == "" {
			continue
		}
		query = query.Preload(preload)
	}
	return query
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return errs.ErrInvalidInput
	}
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository[T]) GetByID(ctx context.Context, id int64, preloads ...string) (*T, error) {
	if id < 1 {
		return nil, errs.ErrInvalidInput
	}

	result := new(T)
	query := r.withPreloads(r.db.WithContext(ctx), preloads)
	if err := query.First(result, id).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *baseRepository[T]) FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error) {
	result := new(T)
	if err := r.db.WithContext(ctx).Where(condition, args...).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *baseRepository[T]) List(ctx context.Context, opts ListOptions, preloads ...string) ([]T, int64, error) {
	normalized, err := normalizeListOptions(opts)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(new(T))
	query, err = applySearchAndSort(query, normalized)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = r.withPreloads(query, preloads)
	offset := (normalized.Page - 1) * normalized.Limit
	query = query.Offset(offset).Limit(normalized.Limit)

	var items []T
	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *baseRepository[T]) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id < 1 || len(updates) == 0 {
		return errs.ErrInvalidInput
	}

	entity := new(T)
	if err := r.db.WithContext(ctx).First(entity, id).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(entity).Updates(updates).Error
}

func (r *baseRepository[T]) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return errs.ErrInvalidInput
	}

	tx := r.db.WithContext(ctx).Delete(new(T), id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func validateLookupValue(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errs.ErrInvalidInput
	}
	return value, nil
}
