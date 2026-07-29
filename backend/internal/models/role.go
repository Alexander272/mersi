package models

import (
	"time"

	"github.com/lib/pq"
)

type Role struct {
	ID          string `json:"id" db:"id"`
	Slug        string `json:"slug" db:"slug"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Level       int    `json:"level" db:"level"`
	IsActive    bool   `json:"isActive" db:"is_active"`
	IsSystem    bool   `json:"isSystem" db:"is_system"`
	IsEditable  bool   `json:"isEditable" db:"is_editable"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type RoleShort struct {
	ID   string `json:"id" db:"id"`
	Slug string `json:"slug" db:"slug"`
	Name string `json:"name" db:"name"`
}

type GetRoleDTO struct {
	ID   string `json:"id" db:"id"`
	Slug string `json:"slug" db:"slug"`
}

type RoleDTO struct {
	ID          string   `json:"id" db:"id"`
	Slug        string   `json:"slug" db:"slug" binding:"required"`
	Name        string   `json:"name" db:"name" binding:"required"`
	Description string   `json:"description" db:"description"`
	Level       int      `json:"level" db:"level"`
	IsSystem    bool     `json:"isSystem" db:"is_system"`
	Permissions []string `json:"permissions" db:"permissions"`
	Inherits    []string `json:"inherits" db:"inherits"`
}

type DeleteRoleDTO struct {
	ID string `json:"id" db:"id"`
}

type RolePermissionDTO struct {
	RoleId       string `json:"roleId" db:"role_id"`
	PermissionId string `json:"permissionId" db:"permission_id"`
}

type RoleWithPerms struct {
	Role
	Inherits []string                  `json:"inherits"`
	Perms    []*RolePermissionsGrouped `json:"perms"`
}

// Legacy types kept for backward compatibility during migration

type RoleLegacy struct {
	ID    string   `json:"id" db:"id"`
	Name  string   `json:"name" db:"name"`
	Rules []string `json:"rules"`
}

type RoleWithRuleDTO struct {
	ID      string         `json:"id" db:"id"`
	Name    string         `json:"name" db:"name"`
	Extends pq.StringArray `db:"extends"`
	Rules   pq.StringArray `db:"rules"`
}

type RoleFull struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Level       int      `json:"level" db:"level"`
	Extends     []string `json:"extends" db:"extends"`
	Description string   `json:"description" db:"description"`
}

type RoleFullDTO struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Level       int            `json:"level" db:"level"`
	Extends     pq.StringArray `json:"extends" db:"extends"`
	Description string         `json:"description" db:"description"`
}

type RoleWithRealm struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Level       int    `json:"level" db:"level"`
	Description string `json:"description" db:"description"`
	RealmId     string `json:"realmId" db:"realm_id"`
}

type RoleWithApi struct{}

type GetRolesDTO struct{}

type GetRoleByRealmDTO struct {
	RealmID string `json:"realmId" binding:"required"`
	UserID  string `json:"userId"`
}

type RoleWithStats struct {
	Role
	Children   []string          `json:"children"`
	UserCount  int               `json:"userCount"`
	PermsCount PermsWithCount    `json:"permsCount"`
}
