package repository

import (
	"context"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	"gorm.io/gorm"
)

type OAuthProviderRepository interface {
	CRUDRepository[models.OAuthProvider]
	GetByProviderName(ctx context.Context, providerName string) (*models.OAuthProvider, error)
}

type UserRepository interface {
	CRUDRepository[models.User]
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
}

type UserOAuthProviderRepository interface {
	CRUDRepository[models.UserOAuthProvider]
}

type RoleRepository interface {
	CRUDRepository[models.Role]
	GetByName(ctx context.Context, roleName string) (*models.Role, error)
}

type UserRoleRepository interface {
	CRUDRepository[models.UserRole]
}

type PermissionRepository interface {
	CRUDRepository[models.Permission]
	GetByType(ctx context.Context, permissionType string) (*models.Permission, error)
}

type RolePermissionRepository interface {
	CRUDRepository[models.RolePermission]
}

type oauthProviderRepository struct {
	*baseRepository[models.OAuthProvider]
}

type userRepository struct {
	*baseRepository[models.User]
}

type userOAuthProviderRepository struct {
	*baseRepository[models.UserOAuthProvider]
}

type roleRepository struct {
	*baseRepository[models.Role]
}

type userRoleRepository struct {
	*baseRepository[models.UserRole]
}

type permissionRepository struct {
	*baseRepository[models.Permission]
}

type rolePermissionRepository struct {
	*baseRepository[models.RolePermission]
}

func NewOAuthProviderRepository(db *gorm.DB) OAuthProviderRepository {
	return &oauthProviderRepository{baseRepository: newBaseRepository[models.OAuthProvider](db)}
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{baseRepository: newBaseRepository[models.User](db)}
}

func NewUserOAuthProviderRepository(db *gorm.DB) UserOAuthProviderRepository {
	return &userOAuthProviderRepository{baseRepository: newBaseRepository[models.UserOAuthProvider](db)}
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{baseRepository: newBaseRepository[models.Role](db)}
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{baseRepository: newBaseRepository[models.UserRole](db)}
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{baseRepository: newBaseRepository[models.Permission](db)}
}

func NewRolePermissionRepository(db *gorm.DB) RolePermissionRepository {
	return &rolePermissionRepository{baseRepository: newBaseRepository[models.RolePermission](db)}
}

func (r *oauthProviderRepository) GetByProviderName(ctx context.Context, providerName string) (*models.OAuthProvider, error) {
	value, err := validateLookupValue(providerName)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "provider_name = ?", value)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	value, err := validateLookupValue(email)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "email = ?", value)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	value, err := validateLookupValue(username)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "username = ?", value)
}

func (r *userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	value, err := validateLookupValue(phoneNumber)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "phone_number = ?", value)
}

func (r *roleRepository) GetByName(ctx context.Context, roleName string) (*models.Role, error) {
	value, err := validateLookupValue(roleName)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "role_name = ?", value)
}

func (r *permissionRepository) GetByType(ctx context.Context, permissionType string) (*models.Permission, error) {
	value, err := validateLookupValue(permissionType)
	if err != nil {
		return nil, err
	}

	return r.FindOne(ctx, "permission_type = ?", value)
}
