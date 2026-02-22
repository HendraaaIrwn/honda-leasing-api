package dto

import "time"

type OAuthProviderDTO struct {
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	IssuerURL    string `json:"issuer_url"`
	Active       bool   `json:"active"`
}

type UserDTO struct {
	UserID         int64      `json:"user_id"`
	Username       string     `json:"username"`
	PhoneNumber    string     `json:"phone_number"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Password       string     `json:"password"`
	PinKey         string     `json:"pin_key"`
	IsActive       bool       `json:"is_active"`
	LastLogin      *time.Time `json:"last_login"`
	FailedAttempts int16      `json:"failed_attempts"`
	LockedUntil    *time.Time `json:"locked_until"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UserOAuthProviderDTO struct {
	UserOAuthID  int64      `json:"user_oauth_id"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UserID       int64      `json:"user_id"`
	ProviderID   int64      `json:"provider_id"`
}

type RoleDTO struct {
	RoleID      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type UserRoleDTO struct {
	UserRoleID int64     `json:"user_role_id"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy int64     `json:"assigned_by"`
	UserID     int64     `json:"user_id"`
	RoleID     int64     `json:"role_id"`
}

type PermissionDTO struct {
	PermissionID   int64  `json:"permission_id"`
	PermissionType string `json:"permission_type"`
	Description    string `json:"description"`
}

type RolePermissionDTO struct {
	RolePermissionID int64 `json:"role_permission_id"`
	RoleID           int64 `json:"role_id"`
	PermissionID     int64 `json:"permission_id"`
}
