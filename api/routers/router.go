package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterERDRouters registers all routers grouped by ERD schema.
func RegisterERDRouters(engine *gin.Engine, basePath string, h *handler.Handlers) {
	root := engine.Group(basePath)

	RegisterMSTRoutes(root.Group("/mst"), h.MST)
	RegisterAccountRoutes(root.Group("/account"), h.Account)
	RegisterDealerRoutes(root.Group("/dealer"), h.Dealer)
	RegisterLeasingRoutes(root.Group("/leasing"), h.Leasing)
	RegisterPaymentRoutes(root.Group("/payment"), h.Payment)
}
