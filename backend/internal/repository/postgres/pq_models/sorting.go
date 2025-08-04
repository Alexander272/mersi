package pq_models

type Sorting struct {
	Id        string `json:"id" db:"id"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
	OrderType string `json:"orderType" db:"order_type"`
	Count     int    `json:"count" db:"count"`
}
