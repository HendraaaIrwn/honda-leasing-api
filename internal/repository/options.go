package repository

import (
	"fmt"
	"regexp"
	"strings"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"gorm.io/gorm"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// ListOptions contains common query options for list endpoints.
type ListOptions struct {
	Page         int
	Limit        int
	SortBy       string
	SortOrder    string
	Search       string
	SearchFields []string

	// AllowedSortFields is optional whitelist for sort_by.
	AllowedSortFields []string
}

var sortFieldPattern = regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`)

func normalizeListOptions(opts ListOptions) (ListOptions, error) {
	if opts.Page == 0 {
		opts.Page = DefaultPage
	}
	if opts.Limit == 0 {
		opts.Limit = DefaultLimit
	}

	if opts.Page < 1 || opts.Limit < 1 {
		return opts, errs.ErrInvalidPagination
	}
	if opts.Limit > MaxLimit {
		opts.Limit = MaxLimit
	}

	if opts.SortBy != "" {
		opts.SortBy = strings.TrimSpace(opts.SortBy)
		if !sortFieldPattern.MatchString(opts.SortBy) {
			return opts, errs.ErrInvalidSort
		}
		if len(opts.AllowedSortFields) > 0 {
			allowed := false
			for _, field := range opts.AllowedSortFields {
				if strings.EqualFold(strings.TrimSpace(field), opts.SortBy) {
					allowed = true
					break
				}
			}
			if !allowed {
				return opts, errs.ErrInvalidSort
			}
		}

		order := strings.ToUpper(strings.TrimSpace(opts.SortOrder))
		if order == "" {
			order = "ASC"
		}
		if order != "ASC" && order != "DESC" {
			return opts, errs.ErrInvalidSort
		}
		opts.SortOrder = order
	}

	opts.Search = strings.TrimSpace(opts.Search)
	if opts.Search != "" && len(opts.SearchFields) == 0 {
		return opts, errs.ErrInvalidSearch
	}

	return opts, nil
}

func applySearchAndSort(query *gorm.DB, opts ListOptions) (*gorm.DB, error) {
	if opts.Search != "" {
		clauses := make([]string, 0, len(opts.SearchFields))
		args := make([]interface{}, 0, len(opts.SearchFields))

		for _, field := range opts.SearchFields {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}
			clauses = append(clauses, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+opts.Search+"%")
		}

		if len(clauses) == 0 {
			return nil, errs.ErrInvalidSearch
		}

		query = query.Where(strings.Join(clauses, " OR "), args...)
	}

	if opts.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", opts.SortBy, opts.SortOrder))
	}

	return query, nil
}
