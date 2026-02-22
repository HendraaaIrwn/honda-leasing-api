package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterDealerRoutes registers routes for schema dealer.*
func RegisterDealerRoutes(group *gin.RouterGroup, h handler.DealerHandlers) {
	registerCRUDRoutes(group, "/motor_types", h.MotorType)
	registerCRUDRoutes(group, "/motors", h.Motor)
	registerCRUDRoutes(group, "/motor_assets", h.MotorAsset)
	registerCRUDRoutes(group, "/customer", h.Customer)
}
