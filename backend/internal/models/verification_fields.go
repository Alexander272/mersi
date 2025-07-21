package models

type VerificationField struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Field     string `json:"field" db:"field"`
	Label     string `json:"label" db:"label"`
	Type      string `json:"type" db:"type"`
	Position  int    `json:"position" db:"position"`
	Group     string `json:"group" db:"group"`
}

type GetVerFieldsDTO struct {
	SectionId string `json:"sectionId" db:"section_id"`
	Group     string `json:"group" db:"group"`
}

type VerificationFieldDTO struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Field     string `json:"field" db:"field"`
	Label     string `json:"label" db:"label"`
	Type      string `json:"type" db:"type"`
	Position  int    `json:"position" db:"position"`
	Group     string `json:"group" db:"group"`
}

type DeleteVerFieldDTO struct {
	Id string `json:"id" db:"id"`
}
