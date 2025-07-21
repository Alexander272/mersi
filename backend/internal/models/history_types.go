package models

import "time"

type HistoryType struct {
	Id       string    `json:"id" db:"id"`
	Group    string    `json:"group" db:"group"`
	Label    string    `json:"label" db:"label"`
	Position int       `json:"position" db:"position"`
	Created  time.Time `json:"created" db:"created_at"`
}

type GetHistoryTypesDTO struct {
	SectionId string `json:"sectionId" db:"section_id"`
}

type HistoryTypeDTO struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Group     string `json:"group" db:"group"`
	Label     string `json:"label" db:"label"`
	Position  int    `json:"position" db:"position" binding:"min=0"`
}

type DeleteHistoryTypeDTO struct {
	Id string `json:"id" db:"id"`
}
