package models

type VerificationField struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Field     string `json:"field" db:"field"`
	Label     string `json:"label" db:"label"`
	Type      string `json:"type" db:"type"`
	Position  int    `json:"position" db:"position"`
}

type GetVerFieldsDTO struct {
	SectionId string `json:"sectionId" db:"section_id"`
}

type VerificationFieldDTO struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Field     string `json:"field" db:"field"`
	Label     string `json:"label" db:"label"`
	Type      string `json:"type" db:"type"`
	Position  int    `json:"position" db:"position"`
}

type DeleteVerFieldDTO struct {
	Id string `json:"id" db:"id"`
}
