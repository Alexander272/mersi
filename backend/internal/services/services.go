package services

import (
	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/Alexander272/mersi/backend/pkg/mattermost"
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
	Location
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
	Filters
	Sorting
	File
	Notification
	Scheduler
	Department
	Employee
	Channel
	Responsible
	Most *most.MostService
	Export
}

type Deps struct {
	Repo          *repository.Repository
	Keycloak      *auth.KeycloakClient
	MostClient    *mattermost.Client
	BotName       string
	CheckUsedConf config.UsedConfig
	// BotUrl   string
}

func NewServices(deps *Deps) *Services {
	role := NewRoleService(deps.Repo.Role)
	ruleItem := NewRuleItemService(deps.Repo.RuleItem)
	rule := NewRuleService(deps.Repo.Rule, ruleItem)

	user := NewUserService(&UsersDeps{Repo: deps.Repo.Users, Keycloak: deps.Keycloak, Role: role})
	session := NewSessionService(deps.Keycloak, user)
	permission := NewPermissionService("configs/privacy.conf", rule, role)

	realm := NewRealmService(deps.Repo.Realm, user)
	accesses := NewAccessesService(deps.Repo.Accesses)
	section := NewSectionService(deps.Repo.Section)
	columns := NewColumnsService(deps.Repo.Columns)
	createForm := NewCreateFormService(deps.Repo.CreateForm)

	document := NewDocumentService(deps.Repo.Document)
	verificationDoc := NewVerificationDocService(deps.Repo.VerificationDoc)
	instrument := NewInstrumentService(deps.Repo.Instrument, document)
	verification := NewVerificationService(deps.Repo.Verification, verificationDoc, instrument, document)

	verificationFields := NewVerificationFieldService(deps.Repo.VerificationFields)
	contextMenu := NewContextService(deps.Repo.ContextMenu, role)
	customContext := NewCustomContextService(deps.Repo.CustomContextMenu)
	toolsMenu := NewToolsMenuService(deps.Repo.ToolsMenu, customContext, role)
	repair := NewRepairService(deps.Repo.Repair, instrument)
	preservation := NewPreservationService(deps.Repo.Preservation)
	transferToSave := NewTransferToSaveService(deps.Repo.TransferToSave)
	transferToDep := NewTransferToDepService(deps.Repo.TransferToDepartment, instrument, document)
	writeOff := NewWriteOffService(deps.Repo.WriteOff, instrument, document)
	historyType := NewHistoryTypeService(deps.Repo.HistoryType)

	filters := NewFilterService(deps.Repo.Filters)
	sorting := NewSortingService(deps.Repo.Sorting)

	responsible := NewResponsibleService(deps.Repo.Responsible)

	//TODO надо бы подумать как избавиться от этой кольцевой зависимости
	si := NewSiService(&SiDeps{
		Repo:         deps.Repo.SI,
		Instrument:   instrument,
		Verification: verification,
		Location:     NewLocationService(&LocationDeps{Repo: deps.Repo.Location, Responsible: responsible}),
	})

	file := NewFileService()
	most := most.NewMostService(most.MostDeps{Client: deps.MostClient})
	notification := NewNotificationService(&NotificationDeps{
		SI:      si,
		File:    file,
		Section: section,
		Most:    most,
		Conf:    deps.CheckUsedConf,
	})
	export := NewExportService(&ExportDeps{
		File:    file,
		SI:      si,
		Columns: columns,
	})
	location := NewLocationService(&LocationDeps{
		Repo:         deps.Repo.Location,
		Responsible:  responsible,
		Notification: notification,
		Most:         most,
	})
	department := NewDepartmentService(deps.Repo.Department, location)
	employee := NewEmployeeService(deps.Repo.Employee, location)
	channel := NewChannelService(deps.Repo.Channel)

	scheduler := NewSchedulerService(&SchedulerDeps{
		Notification: notification,
		User:         user,
		Location:     location,
		Documents:    document,
	})

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
		Filters:              filters,
		Sorting:              sorting,
		Department:           department,
		Employee:             employee,
		Channel:              channel,
		Responsible:          responsible,
		Location:             location,

		File:         file,
		Notification: notification,
		Scheduler:    scheduler,
		Most:         most,
		Export:       export,
	}
}
