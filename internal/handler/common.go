package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func parseIDParam(c *gin.Context, key string) (int64, error) {
	value := strings.TrimSpace(c.Param(key))
	if value == "" {
		return 0, errs.ErrInvalidInput
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, errs.ErrInvalidInput
	}

	return id, nil
}

func parseListRequest(c *gin.Context) (repository.ListOptions, []string, error) {
	page, err := parsePositiveIntQuery(c, "page", repository.DefaultPage)
	if err != nil {
		return repository.ListOptions{}, nil, err
	}

	limit, err := parsePositiveIntQuery(c, "limit", repository.DefaultLimit)
	if err != nil {
		return repository.ListOptions{}, nil, err
	}

	opts := repository.ListOptions{
		Page:      page,
		Limit:     limit,
		SortBy:    strings.TrimSpace(c.Query("sort_by")),
		SortOrder: strings.TrimSpace(c.Query("sort_order")),
		Search:    strings.TrimSpace(c.Query("search")),
	}

	preloads := parsePreloads(c.Query("preload"))
	return opts, preloads, nil
}

func parsePositiveIntQuery(c *gin.Context, key string, fallback int) (int, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, errs.ErrInvalidPagination
	}

	return parsed, nil
}

func parsePreloads(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	raw := strings.Split(value, ",")
	preloads := make([]string, 0, len(raw))
	for _, part := range raw {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		preloads = append(preloads, part)
	}

	if len(preloads) == 0 {
		return nil
	}

	return preloads
}

func paginationMeta(opts repository.ListOptions, total int64) response.PaginationMeta {
	totalPages := 0
	if opts.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(opts.Limit)))
	}

	return response.PaginationMeta{
		Page:       opts.Page,
		Limit:      opts.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

func respondError(c *gin.Context, err error) {
	if err == nil {
		response.InternalServerError(c, "unexpected empty error", nil)
		return
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.NotFound(c, "data not found", err.Error())
	case errors.Is(err, errs.ErrInvalidInput),
		errors.Is(err, errs.ErrInvalidPagination),
		errors.Is(err, errs.ErrInvalidSort),
		errors.Is(err, errs.ErrInvalidSearch),
		errors.Is(err, errs.ErrInvalidPassword),
		errors.Is(err, errs.ErrInvalidDecision),
		errors.Is(err, errs.ErrInvalidStatusTransition),
		errors.Is(err, errs.ErrContractNotDraft),
		errors.Is(err, errs.ErrContractNotApproved),
		errors.Is(err, errs.ErrDPOutOfRange),
		errors.Is(err, errs.ErrInvalidPaymentAmount):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, errs.ErrInvalidEmail), isDuplicateKeyError(err):
		response.Conflict(c, "duplicate data", err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", err.Error())
	}
}

func isDuplicateKeyError(err error) bool {
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}

	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique")
}
