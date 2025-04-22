package models

import "time"

type Accesses struct {
	ID      string    `json:"id" db:"id"`
	RealmID string    `json:"realmId" db:"realm_id"`
	User    *UserData `json:"user"`
	Role    *Role     `json:"role"`
	Created time.Time `json:"created" db:"created_at"`
	// UserId  string `json:"userId" db:"user_id"`
	// RoleId  string `json:"roleId" db:"role_id"`
}

type GetAccessesDTO struct {
	RealmID string `json:"id"`
}

type AccessesDTO struct {
	ID      string `json:"id" db:"id"`
	RealmID string `json:"realmId" db:"realm_id" binding:"required"`
	UserID  string `json:"userId" db:"user_id" binding:"required"`
	RoleID  string `json:"roleId" db:"role_id" binding:"required"`
}

type DeleteAccessesDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}
