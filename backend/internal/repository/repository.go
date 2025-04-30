package repository

import (
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Realm interface {
	postgres.Realm
}
type Section interface {
	postgres.Section
}
type Columns interface {
	postgres.Columns
}
type CreateForm interface {
	postgres.CreateForm
}

type Repository struct {
	Realm
	Section
	Columns
	CreateForm
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Realm:      postgres.NewRealmRepo(db),
		Section:    postgres.NewSectionRepo(db),
		Columns:    postgres.NewColumnRepo(db),
		CreateForm: postgres.NewCreateFormRepo(db),
	}
}
