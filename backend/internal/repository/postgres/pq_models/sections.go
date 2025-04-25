package pq_models

import "time"

type Section struct {
	ID         string    `json:"id" db:"id"`
	RealmID    string    `json:"realmId" db:"realm_id"`
	Realm      string    `db:"realm"`
	RealmTitle string    `db:"title"`
	Name       string    `json:"name" db:"name"`
	Position   int       `json:"position" db:"position" binding:"required"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}
