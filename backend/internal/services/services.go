package services

import (
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/pkg/auth"
)

type Services struct {
	RuleItem
	Rule
	Role
	User
	Session
	Permission

	Realm
	Section
	Columns
	CreateForm
	Instrument
	Document
	VerificationDoc
	Verification
	SI
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
	instrument := NewInstrumentService(deps.Repo.Instrument)
	document := NewDocumentService(deps.Repo.Document)
	verificationDoc := NewVerificationDocService(deps.Repo.VerificationDoc)
	verification := NewVerificationService(deps.Repo.Verification, verificationDoc)

	si := NewSiService(&SiDeps{Repo: deps.Repo.SI, Instrument: instrument, Verification: verification})

	role := NewRoleService(deps.Repo.Role)
	ruleItem := NewRuleItemService(deps.Repo.RuleItem)
	rule := NewRuleService(deps.Repo.Rule, ruleItem)

	user := NewUserService(role)
	session := NewSessionService(deps.Keycloak, user)
	permission := NewPermissionService("configs/privacy.conf", rule, role)

	return &Services{
		Realm:           realm,
		Section:         section,
		Columns:         columns,
		CreateForm:      createForm,
		Instrument:      instrument,
		Document:        document,
		VerificationDoc: verificationDoc,
		Verification:    verification,
		SI:              si,

		Role:     role,
		RuleItem: ruleItem,
		Rule:     rule,

		User:       user,
		Session:    session,
		Permission: permission,
	}
}
