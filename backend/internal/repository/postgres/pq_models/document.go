package pq_models

type Document struct {
	Id           string      `json:"id" db:"id"`
	Label        string      `json:"label" db:"label"`
	Size         int64       `json:"size" db:"size"`
	Path         string      `json:"path" db:"path"`
	DocumentType string      `json:"type" db:"type"`
	Group        string      `json:"group" db:"belongs"`
	InstrumentId interface{} `json:"-" db:"instrument_id"`
	UserId       string      `json:"userId" db:"user_id"`
}
