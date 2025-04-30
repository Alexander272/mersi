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
	ID        string    `json:"id" db:"id"`
	SectionID string    `json:"sectionId" db:"section_id"`
	Step      int       `json:"step" db:"step"`
	StepName  string    `json:"stepName" db:"step_name"`
	Field     string    `json:"field" db:"field"`
	Type      string    `json:"type" db:"type"`
	Position  int       `json:"position" db:"position"`
	Created   time.Time `json:"created" db:"created_at"`
}

type CreateFormFieldDTO struct {
	ID        string `json:"id" db:"id"`
	SectionID string `json:"sectionId" db:"section_id"`
	Step      int    `json:"step" db:"step"`
	StepName  string `json:"stepName" db:"step_name"`
	Field     string `json:"field" db:"field"`
	Type      string `json:"type" db:"type"`
	Position  int    `json:"position" db:"position"`
}

type DeleteCreateFormFieldDTO struct {
	ID string `json:"id" db:"id"`
}
