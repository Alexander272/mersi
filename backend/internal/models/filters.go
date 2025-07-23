package models

type SavedFilter struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	CompareType string `json:"compareType" db:"compare_type"`
	Value       string `json:"value" db:"value"`
}

type GetSavedFiltersDTO struct {
	UserId    string
	SectionId string
}

type SavedFilterDTO struct {
	Id          string `json:"id" db:"id"`
	UserId      string `json:"userId" db:"user_id"`
	SectionId   string `json:"sectionId" db:"section_id"`
	Name        string `json:"name" db:"name"`
	CompareType string `json:"compareType" db:"compare_type"`
	Value       string `json:"value" db:"value"`
}

type DeleteSavedFiltersDTO struct {
	UserId    string
	SectionId string
}
