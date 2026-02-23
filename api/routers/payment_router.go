package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers routes for schema payment.*
func RegisterPaymentRoutes(group *gin.RouterGroup, h handler.PaymentHandlers) {
	registerCRUDRoutes(group, "/payment_schedule", h.PaymentSchedule)
	registerCRUDRoutes(group, "/payments", h.Payment)
}
