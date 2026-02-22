package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

const (
	maxRegionNameLength  = 85
	maxPostalCodeLength  = 10
	maxTemplateName      = 85
	defaultLookupRoleTag = "role"
)

func fetchAllRecords[T any](ctx context.Context, repo repository.CRUDRepository[T]) ([]T, error) {
	result := make([]T, 0)
	page := 1
	for {
		items, total, err := repo.List(ctx, repository.ListOptions{
			Page:  page,
			Limit: repository.MaxLimit,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		if len(items) == 0 || int64(len(result)) >= total {
			break
		}
		page++
	}
	return result, nil
}

func validatePositiveID(label string, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%s ID cannot be empty", strings.TrimSpace(label))
	}
	return nil
}

func validateName(label, value string, maxLength int) error {
	label = strings.TrimSpace(label)
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s name cannot be empty", label)
	}
	if len(value) > maxLength {
		return fmt.Errorf("%s name cannot exceed %d characters", label, maxLength)
	}
	return nil
}

func validatePostalCode(postalCode string) error {
	postalCode = strings.TrimSpace(postalCode)
	if postalCode == "" {
		return nil
	}
	if len(postalCode) > maxPostalCodeLength {
		return fmt.Errorf("postal code cannot exceed %d characters", maxPostalCodeLength)
	}
	return nil
}
