package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterMSTRoutes registers routes for schema mst.*
func RegisterMSTRoutes(group *gin.RouterGroup, h handler.MSTHandlers) {
	registerCRUDRoutes(group, "/province", h.Province)
	registerCRUDRoutes(group, "/kabupaten", h.Kabupaten)
	registerCRUDRoutes(group, "/kecamatan", h.Kecamatan)
	registerCRUDRoutes(group, "/kelurahan", h.Kelurahan)
	registerCRUDRoutes(group, "/locations", h.Location)
	registerCRUDRoutes(group, "/template_tasks", h.TemplateTask)
	registerCRUDRoutes(group, "/template_task_attributes", h.TemplateTaskAttribute)
}
