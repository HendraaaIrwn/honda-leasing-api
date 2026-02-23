package handler

import (
	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
)

type AccountHandlers struct {
	OAuthProvider     ResourceHandler
	User              ResourceHandler
	UserOAuthProvider ResourceHandler
	Role              ResourceHandler
	UserRole          ResourceHandler
	Permission        ResourceHandler
	RolePermission    ResourceHandler
}

func NewAccountHandlers(s services.AccountServices) AccountHandlers {
	return AccountHandlers{
		OAuthProvider:     NewCRUDHandler[models.OAuthProvider]("oauth provider", s.OAuthProvider),
		User:              NewCRUDHandler[models.User]("user", s.User),
		UserOAuthProvider: NewCRUDHandler[models.UserOAuthProvider]("user oauth provider", s.UserOAuthProvider),
		Role:              NewCRUDHandler[models.Role]("role", s.Role),
		UserRole:          NewCRUDHandler[models.UserRole]("user role", s.UserRole),
		Permission:        NewCRUDHandler[models.Permission]("permission", s.Permission),
		RolePermission:    NewCRUDHandler[models.RolePermission]("role permission", s.RolePermission),
	}
}
