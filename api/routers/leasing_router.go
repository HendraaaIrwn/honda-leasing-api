package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterLeasingRoutes registers routes for schema leasing.*
func RegisterLeasingRoutes(group *gin.RouterGroup, h handler.LeasingHandlers) {
	registerCRUDRoutes(group, "/leasing_product", h.LeasingProduct)
	registerCRUDRoutes(group, "/leasing_contract", h.LeasingContract)
	registerCRUDRoutes(group, "/leasing_tasks", h.LeasingTask)
	registerCRUDRoutes(group, "/leasing_tasks_attributes", h.LeasingTaskAttribute)
	registerCRUDRoutes(group, "/leasing_contract_documents", h.LeasingContractDocument)

	if h.Workflow != nil {
		h.Workflow.RegisterRoutes(group)
	}
}
