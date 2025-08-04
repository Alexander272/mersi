package models

type Sorting struct {
	Id        string `json:"id"`
	SectionId string `json:"sectionId"`
	Name      string `json:"name"`
	OrderType string `json:"orderType"`
}

type SortingMap map[string]string

type GetSortingDTO struct {
	UserId    string `json:"userId" db:"sso_id"`
	SectionId string `json:"sectionId" db:"section_id"`
}

type SortingDTO struct {
	Id        string `json:"id" db:"id"`
	UserId    string `json:"userId" db:"sso_id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
	OrderType string `json:"orderType" db:"order_type"`
}

type DeleteSortingDTO struct {
	UserId    string `json:"userId" db:"sso_id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
}
