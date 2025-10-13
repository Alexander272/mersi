package models

type SiStatus struct {
	Id        string `db:"id" json:"id"`
	SectionId string `db:"section_id" json:"sectionId"`
	Position  int    `db:"position" json:"position"`
	Value     string `db:"value" json:"value"`
	Label     string `db:"label" json:"label"`
}

type GetSiStatusDTO struct {
	SectionId string `db:"section_id" json:"sectionId"`
}

type SiStatusDTO struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Position  int    `json:"position" db:"position"`
	Value     string `json:"value" db:"value"`
	Label     string `json:"label" db:"label"`
}

type DeleteSiStatusDTO struct {
	Id string `json:"id" db:"id"`
}
