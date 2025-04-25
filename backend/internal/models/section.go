package models

import "time"

type Section struct {
	ID        string    `json:"id" db:"id"`
	RealmID   string    `json:"realmId" db:"realm_id"`
	Name      string    `json:"name" db:"name"`
	Position  int       `json:"position" db:"position" binding:"required"`
	CreatedAt time.Time `json:"created" db:"created_at"`
}

type GroupedSections struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Realm    string     `json:"realm"`
	Sections []*Section `json:"sections"`
}

type GetGroupedSectionDTO struct{}

type GetSectionsDTO struct {
	RealmID string `json:"realmId"`
}

type SectionDTO struct {
	ID          string `json:"id" db:"id" binding:"required"`
	Name        string `json:"name" db:"name" binding:"required"`
	Position    int    `json:"position" db:"position" binding:"required"`
	MaxPosition int    `json:"maxPosition" db:"max_position"`
	RealmID     string `json:"realmId" db:"realm_id" binding:"required"`
}

type DeleteSectionDTO struct {
	ID string `json:"id" db:"id" binding:"required"`
}
