package models

import "time"

type Realm struct {
	ID                   string    `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	Realm                string    `json:"realm" db:"realm"`
	IsActive             bool      `json:"isActive" db:"is_active"`
	ReserveChannel       string    `json:"reserveChannel" db:"reserve_channel"`
	ExpirationNotice     bool      `json:"expirationNotice" db:"expiration_notice"`
	LocationType         string    `json:"locationType" db:"location_type"`
	NeedConfirmed        bool      `json:"needConfirmed" db:"need_confirmed"`
	NasEmployees         bool      `json:"hasEmployees" db:"has_employees"`
	HasResponsible       bool      `json:"hasResponsible" db:"has_responsible"`
	HasCommissioningCert bool      `json:"hasCommissioningCert" db:"has_commissioning_cert"`
	HasPreservations     bool      `json:"hasPreservations" db:"has_preservations"`
	HasTransfer          bool      `json:"hasTransfer" db:"has_transfer"`
	Created              time.Time `json:"created" db:"created_at"`
}

type GetRealmsDTO struct {
	All bool `json:"all"`
}

type GetRealmByIdDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}

type GetRealmByUserDTO struct {
	UserID string `json:"userId" db:"user_id"`
}

type ChooseRealmDTO struct {
	RealmID string `json:"realmId" db:"realm_id" binding:"required"`
	UserID  string `json:"userId"`
	Role    string `json:"role"`
}

type RealmDTO struct {
	ID                   string `json:"id" db:"id"`
	Name                 string `json:"name" db:"name" binding:"required"`
	Realm                string `json:"realm" db:"realm" binding:"required"`
	IsActive             bool   `json:"isActive" db:"is_active"`
	ReserveChannel       string `json:"reserveChannel" db:"reserve_channel"`
	ExpirationNotice     bool   `json:"expirationNotice" db:"expiration_notice"`
	LocationType         string `json:"locationType" db:"location_type"`
	NeedConfirmed        bool   `json:"needConfirmed" db:"need_confirmed"`
	NasEmployees         bool   `json:"hasEmployees" db:"has_employees"`
	HasResponsible       bool   `json:"hasResponsible" db:"has_responsible"`
	HasCommissioningCert bool   `json:"hasCommissioningCert" db:"has_commissioning_cert"`
	HasPreservations     bool   `json:"hasPreservations" db:"has_preservations"`
	HasTransfer          bool   `json:"hasTransfer" db:"has_transfer"`
}

type DeleteRealmDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}
