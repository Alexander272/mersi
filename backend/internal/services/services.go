package services

import (
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/pkg/auth"
)

type Services struct {
	Realm
	Section
	Columns
	CreateForm
}

type Deps struct {
	Repo     *repository.Repository
	Keycloak *auth.KeycloakClient
	BotUrl   string
}

func NewServices(deps *Deps) *Services {
	realm := NewRealmService(deps.Repo.Realm)
	section := NewSectionService(deps.Repo.Section)
	columns := NewColumnsService(deps.Repo.Columns)
	createForm := NewCreateFormService(deps.Repo.CreateForm)

	return &Services{
		Realm:      realm,
		Section:    section,
		Columns:    columns,
		CreateForm: createForm,
	}
}
