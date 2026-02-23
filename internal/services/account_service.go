package services

import (
	"context"
	"fmt"

	"github.com/HendraaaIrwn/honda-leasing-api/internal/domain/models"
	errs "github.com/HendraaaIrwn/honda-leasing-api/internal/errors"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
)

type OAuthProviderService interface {
	CRUDService[models.OAuthProvider]
	GetByProviderName(ctx context.Context, providerName string) (*models.OAuthProvider, error)
}

type UserService interface {
	CRUDService[models.User]
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
}

type UserOAuthProviderService interface {
	CRUDService[models.UserOAuthProvider]
}

type RoleService interface {
	CRUDService[models.Role]
	GetByName(ctx context.Context, roleName string) (*models.Role, error)
}

type UserRoleService interface {
	CRUDService[models.UserRole]
}

type PermissionService interface {
	CRUDService[models.Permission]
	GetByType(ctx context.Context, permissionType string) (*models.Permission, error)
}

type RolePermissionService interface {
	CRUDService[models.RolePermission]
}

type oauthProviderService struct {
	*baseService[models.OAuthProvider]
	repo repository.OAuthProviderRepository
}

type userService struct {
	*baseService[models.User]
	repo repository.UserRepository
}

type userOAuthProviderService struct {
	*baseService[models.UserOAuthProvider]
}

type roleService struct {
	*baseService[models.Role]
	repo repository.RoleRepository
}

type userRoleService struct {
	*baseService[models.UserRole]
}

type permissionService struct {
	*baseService[models.Permission]
	repo repository.PermissionRepository
}

type rolePermissionService struct {
	*baseService[models.RolePermission]
}

func NewOAuthProviderService(repo repository.OAuthProviderRepository) OAuthProviderService {
	return &oauthProviderService{
		baseService: newBaseService[models.OAuthProvider](repo),
		repo:        repo,
	}
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		baseService: newBaseService[models.User](repo),
		repo:        repo,
	}
}

func NewUserOAuthProviderService(repo repository.UserOAuthProviderRepository) UserOAuthProviderService {
	return &userOAuthProviderService{
		baseService: newBaseService[models.UserOAuthProvider](repo),
	}
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{
		baseService: newBaseService[models.Role](repo),
		repo:        repo,
	}
}

func NewUserRoleService(repo repository.UserRoleRepository) UserRoleService {
	return &userRoleService{
		baseService: newBaseService[models.UserRole](repo),
	}
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{
		baseService: newBaseService[models.Permission](repo),
		repo:        repo,
	}
}

func NewRolePermissionService(repo repository.RolePermissionRepository) RolePermissionService {
	return &rolePermissionService{
		baseService: newBaseService[models.RolePermission](repo),
	}
}

func (s *oauthProviderService) GetByProviderName(ctx context.Context, providerName string) (*models.OAuthProvider, error) {
	return s.repo.GetByProviderName(ctx, providerName)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *userService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	return s.repo.GetByPhoneNumber(ctx, phoneNumber)
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errs.ErrInvalidInput
	}

	hashedPassword, err := hashPasswordWithBcrypt(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %w", err)
	}
	user.Password = hashedPassword

	return s.repo.Create(ctx, user)
}

func (s *userService) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id < 1 || len(updates) == 0 {
		return errs.ErrInvalidInput
	}

	normalized := make(map[string]interface{}, len(updates))
	for key, value := range updates {
		normalized[key] = value
	}

	for _, key := range []string{"password", "Password"} {
		value, ok := normalized[key]
		if !ok {
			continue
		}

		rawPassword, ok := value.(string)
		if !ok {
			return errs.ErrInvalidPassword
		}

		hashedPassword, err := hashPasswordWithBcrypt(rawPassword)
		if err != nil {
			return fmt.Errorf("failed to hash user password: %w", err)
		}

		normalized[key] = hashedPassword
	}

	return s.repo.Update(ctx, id, normalized)
}

func (s *roleService) GetByName(ctx context.Context, roleName string) (*models.Role, error) {
	return s.repo.GetByName(ctx, roleName)
}

func (s *permissionService) GetByType(ctx context.Context, permissionType string) (*models.Permission, error) {
	return s.repo.GetByType(ctx, permissionType)
}
