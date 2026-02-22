package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerCRUDRoutes(group *gin.RouterGroup, resource string, h handler.ResourceHandler) {
	group.GET(resource, h.List)
	group.GET(resource+"/:id", h.GetByID)
	group.POST(resource, h.Create)
	group.PUT(resource+"/:id", h.Update)
	group.DELETE(resource+"/:id", h.Delete)
}
