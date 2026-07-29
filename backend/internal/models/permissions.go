package models

import (
	"time"
)

type Permission struct {
	Id          string    `json:"id" db:"id"`
	Realm       string    `json:"realm" db:"realm"`
	Role        string    `json:"role" db:"role"`
	Object      string    `json:"object" db:"object"`
	Action      string    `json:"action" db:"action"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type GroupedPermission struct {
	Group string        `json:"group"`
	Title string        `json:"title"`
	Items []*Permission `json:"items"`
}

type GetPermsByRoleDTO struct {
	Role  string `json:"role" db:"role"`
	Realm string `json:"realm" db:"realm"`
}

type PermissionDTO struct {
	Id          string `json:"id" db:"id"`
	Object      string `json:"object" db:"object"`
	Action      string `json:"action" db:"action"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type DeletePermissionDTO struct {
	ActorId string `json:"actorId" db:"actor_id"`
	Id      string `json:"id" db:"id"`
}

type RolePermissionItem struct {
	PermissionId string `json:"permissionId" db:"permission_id"`
	Object       string `json:"object" db:"object"`
	Action       string `json:"action" db:"action"`
	IsAssigned   bool   `json:"isAssigned"`
	IsInherited  bool   `json:"isInherited"`
}

type RolePermissionsGrouped struct {
	Group     string                `json:"group"`
	Title     string                `json:"title"`
	Resources []*RolePermissionItem `json:"resources"`
}

type PermsWithCount struct {
	Own       Perm `json:"own"`
	Inherited Perm `json:"inherited"`
	Total     Perm `json:"total"`
}

type Perm struct {
	Items []string `json:"items"`
	Count int      `json:"count"`
}
