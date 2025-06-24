package pq_models

import "github.com/lib/pq"

type RoleFull struct {
	Id          string         `db:"id"`
	Name        string         `db:"name"`
	Level       int            ` db:"level"`
	Extends     pq.StringArray `db:"extends"`
	Description string         `db:"description"`
	RealmId     string         `db:"realm_id"`
}
