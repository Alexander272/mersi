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
	Accesses
	Section
	Columns
	CreateForm
	Instrument
	Document
	VerificationDoc
	Verification
	SI
	ContextMenu
	ToolsMenu
	Repair
	VerificationFields
	Preservation
	TransferToSave
	TransferToDepartment
	WriteOff
	HistoryType
}

type Deps struct {
	Repo     *repository.Repository
	Keycloak *auth.KeycloakClient
	BotUrl   string
}

func NewServices(deps *Deps) *Services {
	role := NewRoleService(deps.Repo.Role)
	ruleItem := NewRuleItemService(deps.Repo.RuleItem)
	rule := NewRuleService(deps.Repo.Rule, ruleItem)

	user := NewUserService(role)
	session := NewSessionService(deps.Keycloak, user)
	permission := NewPermissionService("configs/privacy.conf", rule, role)

	realm := NewRealmService(deps.Repo.Realm)
	accesses := NewAccessesService(deps.Repo.Accesses)
	section := NewSectionService(deps.Repo.Section)
	columns := NewColumnsService(deps.Repo.Columns)
	createForm := NewCreateFormService(deps.Repo.CreateForm)
	instrument := NewInstrumentService(deps.Repo.Instrument)
	document := NewDocumentService(deps.Repo.Document)
	verificationDoc := NewVerificationDocService(deps.Repo.VerificationDoc)
	verification := NewVerificationService(deps.Repo.Verification, verificationDoc)

	si := NewSiService(&SiDeps{Repo: deps.Repo.SI, Instrument: instrument, Verification: verification})

	verificationFields := NewVerificationFieldService(deps.Repo.VerificationFields)
	contextMenu := NewContextService(deps.Repo.ContextMenu, role)
	customContext := NewCustomContextService(deps.Repo.CustomContextMenu)
	toolsMenu := NewToolsMenuService(deps.Repo.ToolsMenu, customContext, role)
	repair := NewRepairService(deps.Repo.Repair)
	preservation := NewPreservationService(deps.Repo.Preservation)
	transferToSave := NewTransferToSaveService(deps.Repo.TransferToSave)
	transferToDep := NewTransferToDepService(deps.Repo.TransferToDepartment, instrument)
	writeOff := NewWriteOffService(deps.Repo.WriteOff)
	historyType := NewHistoryTypeService(deps.Repo.HistoryType)

	return &Services{
		Role:     role,
		RuleItem: ruleItem,
		Rule:     rule,

		User:       user,
		Session:    session,
		Permission: permission,

		Realm:                realm,
		Accesses:             accesses,
		Section:              section,
		Columns:              columns,
		CreateForm:           createForm,
		Instrument:           instrument,
		Document:             document,
		VerificationDoc:      verificationDoc,
		Verification:         verification,
		SI:                   si,
		ContextMenu:          contextMenu,
		ToolsMenu:            toolsMenu,
		Repair:               repair,
		VerificationFields:   verificationFields,
		Preservation:         preservation,
		TransferToSave:       transferToSave,
		TransferToDepartment: transferToDep,
		WriteOff:             writeOff,
		HistoryType:          historyType,
	}
}
