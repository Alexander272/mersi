package models

import "time"

type CreateFormStep struct {
	Step     int                `json:"step"`
	StepName string             `json:"stepName"`
	Fields   []*CreateFormField `json:"fields"`
}

type GetCreateFormDTO struct {
	SectionID string `json:"sectionId"`
}

type CreateFormField struct {
	ID         string    `json:"id" db:"id"`
	SectionID  string    `json:"sectionId" db:"section_id"`
	Step       int       `json:"step" db:"step"`
	StepName   string    `json:"stepName" db:"step_name"`
	Field      string    `json:"field" db:"field"`
	FieldName  string    `json:"fieldName" db:"field_name"`
	Path       string    `json:"path" db:"path"`
	Type       string    `json:"type" db:"type"`
	IsRequired bool      `json:"isRequired" db:"is_required"`
	Position   int       `json:"position" db:"position"`
	Created    time.Time `json:"created" db:"created_at"`
}

type CreateFormFieldDTO struct {
	ID         string `json:"id" db:"id"`
	SectionID  string `json:"sectionId" db:"section_id" binding:"required"`
	Step       int    `json:"step" db:"step" binding:"min=0"`
	StepName   string `json:"stepName" db:"step_name" binding:"required"`
	Field      string `json:"field" db:"field" binding:"required"`
	FieldName  string `json:"fieldName" db:"field_name" binding:"required"`
	Path       string `json:"path" db:"path"`
	Type       string `json:"type" db:"type" binding:"required"`
	IsRequired bool   `json:"isRequired" db:"is_required"`
	Position   int    `json:"position" db:"position" binding:"min=0"`
}

type DeleteCreateFormFieldDTO struct {
	ID string `json:"id" db:"id"`
}
