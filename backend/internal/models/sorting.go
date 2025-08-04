package models

type Sorting struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
	OrderType string `json:"orderType" db:"order_type"`
	Count     int    `json:"count" db:"count"`
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
	Count     int    `json:"count" db:"count"`
}

type DeleteSortingDTO struct {
	UserId    string `json:"userId" db:"sso_id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
}
