package models

import "github.com/lib/pq"

type Role struct {
	ID   string   `json:"id" db:"id"`
	Name string   `json:"name" db:"name"`
	Menu []string `json:"menu"`
}
type RoleWithMenuDTO struct {
	ID      string         `json:"id" db:"id"`
	Name    string         `json:"name" db:"name"`
	Extends pq.StringArray `db:"extends"`
	Menu    pq.StringArray `db:"menu"`
	// Menu    string         `db:"menu"`
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

type RoleDTO struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name" binding:"required"`
	Level       int      `json:"level" db:"level" binding:"required"`
	Extends     []string `json:"extends" db:"extends"`
	Description string   `json:"description" db:"description"`
}
