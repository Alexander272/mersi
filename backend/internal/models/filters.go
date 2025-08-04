package models

type SavedFilter struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"field" db:"name"`
	FieldType   string `json:"fieldType" db:"field_type"`
	CompareType string `json:"compareType" db:"compare_type"`
	Value       string `json:"value" db:"value"`
}

type GetSavedFiltersDTO struct {
	UserId    string
	SectionId string
}

type ChangeFillersDTO struct {
	SectionId string            `json:"section"`
	Filters   []*SavedFilterDTO `json:"filters"`
}

type SavedFilterDTO struct {
	Id          string `json:"id" db:"id"`
	UserId      string `json:"userId" db:"sso_id"`
	SectionId   string `json:"sectionId" db:"section_id"`
	Name        string `json:"field" db:"name"`
	FieldType   string `json:"fieldType" db:"field_type"`
	CompareType string `json:"compareType" db:"compare_type"`
	Value       string `json:"value" db:"value"`
}

type DeleteSavedFiltersDTO struct {
	UserId    string `json:"userId" db:"sso_id"`
	SectionId string `json:"sectionId" db:"section_id"`
}
