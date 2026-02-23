package services

import (
	"context"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

// CRUDService defines common service operations for entities.
type CRUDService[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id int64, preloads ...string) (*T, error)
	FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
	List(ctx context.Context, opts repository.ListOptions, preloads ...string) ([]T, int64, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
}

type baseService[T any] struct {
	repo repository.CRUDRepository[T]
}

func newBaseService[T any](repo repository.CRUDRepository[T]) *baseService[T] {
	return &baseService[T]{repo: repo}
}

func (s *baseService[T]) Create(ctx context.Context, entity *T) error {
	return s.repo.Create(ctx, entity)
}

func (s *baseService[T]) GetByID(ctx context.Context, id int64, preloads ...string) (*T, error) {
	return s.repo.GetByID(ctx, id, preloads...)
}

func (s *baseService[T]) FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error) {
	return s.repo.FindOne(ctx, condition, args...)
}

func (s *baseService[T]) List(ctx context.Context, opts repository.ListOptions, preloads ...string) ([]T, int64, error) {
	return s.repo.List(ctx, opts, preloads...)
}

func (s *baseService[T]) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	return s.repo.Update(ctx, id, updates)
}

func (s *baseService[T]) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
