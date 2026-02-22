package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Error   *ErrorPayload `json:"error,omitempty"`
	Meta    interface{}   `json:"meta,omitempty"`
}

// ErrorPayload describes API error detail.
type ErrorPayload struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationMeta is optional pagination metadata for list endpoints.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// JSON returns a success response with optional data and meta.
func JSON(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error returns a standardized error response.
func Error(c *gin.Context, status int, code, message string, details interface{}) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &ErrorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func OK(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusCreated, message, data, nil)
}

func Paginated(c *gin.Context, message string, data interface{}, meta PaginationMeta) {
	JSON(c, http.StatusOK, message, data, meta)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, details)
}

func Forbidden(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, details)
}

func NotFound(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, details)
}

func Conflict(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusConflict, "CONFLICT", message, details)
}

func UnprocessableEntity(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, details)
}

func InternalServerError(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details)
}
