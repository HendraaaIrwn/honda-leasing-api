package models

import "time"

type OAuthProvider struct {
	ProviderID         int64               `gorm:"column:provider_id;primaryKey;autoIncrement"`
	ProviderName       string              `gorm:"column:provider_name;size:50;not null"`
	ClientID           string              `gorm:"column:client_id;size:255;not null"`
	ClientSecret       string              `gorm:"column:client_secret;size:255;not null"`
	RedirectURI        string              `gorm:"column:redirect_uri;size:255;not null"`
	IssuerURL          string              `gorm:"column:issuer_url;size:255;not null"`
	Active             bool                `gorm:"column:active;not null;default:true"`
	UserOAuthProviders []UserOAuthProvider `gorm:"foreignKey:ProviderID;references:ProviderID"`
}

func (OAuthProvider) TableName() string { return "account.oauth_providers" }

type User struct {
	UserID             int64               `gorm:"column:user_id;primaryKey;autoIncrement"`
	Username           string              `gorm:"column:username;size:50;not null;uniqueIndex"`
	PhoneNumber        string              `gorm:"column:phone_number;size:15;not null;uniqueIndex"`
	Email              string              `gorm:"column:email;size:100;not null;uniqueIndex"`
	FullName           string              `gorm:"column:full_name;size:100;not null"`
	Password           string              `gorm:"column:password;size:255;not null"`
	PinKey             string              `gorm:"column:pin_key;size:255;not null"`
	IsActive           bool                `gorm:"column:is_active;not null;default:true"`
	LastLogin          *time.Time          `gorm:"column:last_login;type:timestamptz"`
	FailedAttempts     int16               `gorm:"column:failed_attempts;not null;default:0"`
	LockedUntil        *time.Time          `gorm:"column:locked_until;type:timestamptz"`
	CreatedAt          time.Time           `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt          time.Time           `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	UserOAuthProviders []UserOAuthProvider `gorm:"foreignKey:UserID;references:UserID"`
	UserRoles          []UserRole          `gorm:"foreignKey:UserID;references:UserID"`
}

func (User) TableName() string { return "account.users" }

type UserOAuthProvider struct {
	UserOAuthID  int64         `gorm:"column:user_oauth_id;primaryKey;autoIncrement"`
	AccessToken  string        `gorm:"column:access_token;type:text;not null"`
	RefreshToken string        `gorm:"column:refresh_token;type:text;not null"`
	ExpiresAt    *time.Time    `gorm:"column:expires_at;type:timestamptz"`
	CreatedAt    time.Time     `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UserID       int64         `gorm:"column:user_id;not null;index"`
	ProviderID   int64         `gorm:"column:provider_id;not null;index"`
	User         User          `gorm:"foreignKey:UserID;references:UserID"`
	Provider     OAuthProvider `gorm:"foreignKey:ProviderID;references:ProviderID"`
}

func (UserOAuthProvider) TableName() string { return "account.user_oauth_provider" }

type Role struct {
	RoleID          int64            `gorm:"column:role_id;primaryKey;autoIncrement"`
	RoleName        string           `gorm:"column:role_name;size:50;not null;uniqueIndex"`
	Description     string           `gorm:"column:description;type:text"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID;references:RoleID"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID;references:RoleID"`
}

func (Role) TableName() string { return "account.roles" }

type UserRole struct {
	UserRoleID int64     `gorm:"column:user_role_id;primaryKey;autoIncrement"`
	AssignedAt time.Time `gorm:"column:assigned_at;type:timestamptz;autoCreateTime"`
	AssignedBy int64     `gorm:"column:assigned_by;not null"`
	UserID     int64     `gorm:"column:user_id;not null;index"`
	RoleID     int64     `gorm:"column:role_id;not null;index"`
	User       User      `gorm:"foreignKey:UserID;references:UserID"`
	Role       Role      `gorm:"foreignKey:RoleID;references:RoleID"`
}

func (UserRole) TableName() string { return "account.user_roles" }

type Permission struct {
	PermissionID    int64            `gorm:"column:permission_id;primaryKey;autoIncrement"`
	PermissionType  string           `gorm:"column:permission_type;size:100;not null;uniqueIndex"`
	Description     string           `gorm:"column:description;type:text"`
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID;references:PermissionID"`
}

func (Permission) TableName() string { return "account.permissions" }

type RolePermission struct {
	RolePermissionID int64      `gorm:"column:role_permission_id;primaryKey;autoIncrement"`
	RoleID           int64      `gorm:"column:role_id;not null;index"`
	PermissionID     int64      `gorm:"column:permission_id;not null;index"`
	Role             Role       `gorm:"foreignKey:RoleID;references:RoleID"`
	Permission       Permission `gorm:"foreignKey:PermissionID;references:PermissionID"`
}

func (RolePermission) TableName() string { return "account.role_permission" }
