package routers

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/gin-gonic/gin"
)

// RegisterAccountRoutes registers routes for schema account.*
func RegisterAccountRoutes(group *gin.RouterGroup, h handler.AccountHandlers) {
	registerCRUDRoutes(group, "/oauth_providers", h.OAuthProvider)
	registerCRUDRoutes(group, "/users", h.User)
	registerCRUDRoutes(group, "/user_oauth_provider", h.UserOAuthProvider)
	registerCRUDRoutes(group, "/roles", h.Role)
	registerCRUDRoutes(group, "/user_roles", h.UserRole)
	registerCRUDRoutes(group, "/permissions", h.Permission)
	registerCRUDRoutes(group, "/role_permission", h.RolePermission)
}
