package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/response"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
	"github.com/gin-gonic/gin"
)

// ResourceHandler defines common REST handlers for a single entity.
type ResourceHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// CRUDHandler is generic REST handler for a single model.
type CRUDHandler[T any] struct {
	name    string
	service services.CRUDService[T]
	mapper  *modelPayloadMapper
}

func NewCRUDHandler[T any](name string, service services.CRUDService[T]) *CRUDHandler[T] {
	return &CRUDHandler[T]{
		name:    strings.TrimSpace(name),
		service: service,
		mapper:  newModelPayloadMapper[T](),
	}
}

func (h *CRUDHandler[T]) List(c *gin.Context) {
	opts, preloads, err := parseListRequest(c)
	if err != nil {
		respondError(c, err)
		return
	}

	items, total, err := h.service.List(c.Request.Context(), opts, preloads...)
	if err != nil {
		respondError(c, err)
		return
	}

	response.Paginated(c, fmt.Sprintf("%s list", h.name), items, paginationMeta(opts, total))
}

func (h *CRUDHandler[T]) GetByID(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	entity, err := h.service.GetByID(c.Request.Context(), id, parsePreloads(c.Query("preload"))...)
	if err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, fmt.Sprintf("%s detail", h.name), entity)
}

func (h *CRUDHandler[T]) Create(c *gin.Context) {
	payload, err := bindPayloadMap(c)
	if err != nil {
		respondError(c, err)
		return
	}

	var entity T
	if err := h.mapper.decodeCreatePayload(payload, &entity); err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.Create(c.Request.Context(), &entity); err != nil {
		respondError(c, err)
		return
	}

	response.Created(c, fmt.Sprintf("%s created", h.name), entity)
}

func (h *CRUDHandler[T]) Update(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	payload, err := bindPayloadMap(c)
	if err != nil {
		respondError(c, err)
		return
	}

	updates, err := h.mapper.buildUpdatePayload(payload)
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.Update(c.Request.Context(), id, updates); err != nil {
		respondError(c, err)
		return
	}

	entity, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.OK(c, fmt.Sprintf("%s updated", h.name), gin.H{"id": id})
		return
	}

	response.OK(c, fmt.Sprintf("%s updated", h.name), entity)
}

func (h *CRUDHandler[T]) Delete(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	response.OK(c, fmt.Sprintf("%s deleted", h.name), gin.H{"id": id})
}

func bindPayloadMap(c *gin.Context) (map[string]interface{}, error) {
	payload := map[string]interface{}{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		return nil, errs.ErrInvalidInput
	}
	if len(payload) == 0 {
		return nil, errs.ErrInvalidInput
	}
	return payload, nil
}

type modelPayloadMapper struct {
	keyToFieldName map[string]string
	keyToColumn    map[string]string
	primaryColumns map[string]struct{}
}

func newModelPayloadMapper[T any]() *modelPayloadMapper {
	mapper := &modelPayloadMapper{
		keyToFieldName: make(map[string]string),
		keyToColumn:    make(map[string]string),
		primaryColumns: make(map[string]struct{}),
	}

	modelType := reflect.TypeOf((*T)(nil)).Elem()
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return mapper
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		column, primary := parseGormField(field.Tag.Get("gorm"), fieldName)

		mapper.registerField(fieldName, fieldName, column)
		mapper.registerField(toSnakeCase(fieldName), fieldName, column)

		if primary {
			mapper.primaryColumns[strings.ToLower(column)] = struct{}{}
		}
	}

	return mapper
}

func (m *modelPayloadMapper) registerField(key, fieldName, column string) {
	if key == "" {
		return
	}

	normalized := normalizePayloadKey(key)
	if normalized == "" {
		return
	}

	if _, exists := m.keyToFieldName[normalized]; !exists {
		m.keyToFieldName[normalized] = fieldName
	}
	if _, exists := m.keyToColumn[normalized]; !exists {
		m.keyToColumn[normalized] = column
	}
}

func (m *modelPayloadMapper) decodeCreatePayload(payload map[string]interface{}, out interface{}) error {
	fieldMap := make(map[string]interface{}, len(payload))
	for rawKey, value := range payload {
		fieldName, ok := m.keyToFieldName[normalizePayloadKey(rawKey)]
		if !ok {
			continue
		}
		fieldMap[fieldName] = value
	}

	if len(fieldMap) == 0 {
		return errs.ErrInvalidInput
	}

	encoded, err := json.Marshal(fieldMap)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(encoded, out); err != nil {
		return errs.ErrInvalidInput
	}

	return nil
}

func (m *modelPayloadMapper) buildUpdatePayload(payload map[string]interface{}) (map[string]interface{}, error) {
	updates := make(map[string]interface{}, len(payload))
	for rawKey, value := range payload {
		normalizedKey := normalizePayloadKey(rawKey)
		column, ok := m.keyToColumn[normalizedKey]
		if !ok {
			continue
		}

		if _, primary := m.primaryColumns[strings.ToLower(column)]; primary {
			continue
		}

		updates[column] = value
	}

	if len(updates) == 0 {
		return nil, errs.ErrInvalidInput
	}

	return updates, nil
}

func parseGormField(tag string, fallbackFieldName string) (column string, primary bool) {
	column = toSnakeCase(fallbackFieldName)
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.EqualFold(part, "primaryKey") {
			primary = true
			continue
		}

		if strings.HasPrefix(strings.ToLower(part), "column:") {
			columnValue := strings.TrimSpace(part[len("column:"):])
			if columnValue != "" {
				column = columnValue
			}
		}
	}

	return column, primary
}

func normalizePayloadKey(key string) string {
	key = strings.TrimSpace(strings.ToLower(key))
	replacer := strings.NewReplacer("_", "", "-", "", " ", "", ".", "")
	return replacer.Replace(key)
}

func toSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	var b strings.Builder
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(input[i-1])
				if prev != '_' && !unicode.IsUpper(prev) {
					b.WriteByte('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
