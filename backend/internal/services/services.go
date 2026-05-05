package services

import (
	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/Alexander272/mersi/backend/pkg/mattermost"
	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
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
	SiStatus
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
	ActivityLog
	Filters
	Sorting
	File
	Notification
	Scheduler
	Department
	DepartmentAccess
	Employee
	Channel
	Responsible
	Most *most.MostService
	Export
	ImportFile
}

type Deps struct {
	Repo          *repository.Repository
	Keycloak      *auth.KeycloakClient
	MostClient    *mattermost.Client
	BotName       string
	CheckUsedConf config.UsedConfig
	Adapter       *sqlxadapter.Adapter
	// BotUrl   string
}

func NewServices(deps *Deps) *Services {
	txManager := NewTransactionManager(deps.Repo.Transaction)

	role := NewRoleService(deps.Repo.Role)
	ruleItem := NewRuleItemService(deps.Repo.RuleItem)
	rule := NewRuleService(deps.Repo.Rule, ruleItem)

	user := NewUserService(&UsersDeps{Repo: deps.Repo.Users, TxManager: txManager, Keycloak: deps.Keycloak, Role: role})
	session := NewSessionService(deps.Keycloak, user)
	realm := NewRealmService(deps.Repo.Realm, user)
	accesses := NewAccessesService(deps.Repo.Accesses)

	permission := NewPermissionService(&PermissionDeps{
		ConfPath: "configs/privacy.conf",
		Adapter:  deps.Adapter,
		Rule:     rule,
		Role:     role,
		Realm:    realm,
		Accesses: accesses,
	})

	section := NewSectionService(deps.Repo.Section)
	columns := NewColumnsService(deps.Repo.Columns)
	createForm := NewCreateFormService(deps.Repo.CreateForm)
	status := NewStatusService(deps.Repo.SiStatus)

	document := NewDocumentService(deps.Repo.Document)
	verificationDoc := NewVerificationDocService(deps.Repo.VerificationDoc)
	activityLog := NewActivityLogService(deps.Repo.ActivityLog)
	instrument := NewInstrumentService(deps.Repo.Instrument, document, activityLog)
	verification := NewVerificationService(&VerificationDeps{
		Repo:       deps.Repo.Verification,
		TxManager:  txManager,
		VerDocs:    verificationDoc,
		Instrument: instrument,
		Docs:       document,
		ActivityLog: activityLog,
	})

	verificationFields := NewVerificationFieldService(deps.Repo.VerificationFields)
	contextMenu := NewContextService(deps.Repo.ContextMenu, role, accesses)
	customContext := NewCustomContextService(deps.Repo.CustomContextMenu)
	toolsMenu := NewToolsMenuService(deps.Repo.ToolsMenu, customContext, role, accesses)
	repair := NewRepairService(deps.Repo.Repair, txManager, instrument, activityLog)
	preservation := NewPreservationService(deps.Repo.Preservation, txManager, instrument, activityLog)
	transferToSave := NewTransferToSaveService(deps.Repo.TransferToSave, txManager, instrument, activityLog)
	transferToDep := NewTransferToDepService(deps.Repo.TransferToDepartment, txManager, instrument, document, activityLog)
	writeOff := NewWriteOffService(deps.Repo.WriteOff, txManager, instrument, document, activityLog)
	historyType := NewHistoryTypeService(deps.Repo.HistoryType)
	// activityLog := NewActivityLogService(deps.Repo.ActivityLog)

	filters := NewFilterService(deps.Repo.Filters)
	sorting := NewSortingService(deps.Repo.Sorting)

	responsible := NewResponsibleService(deps.Repo.Responsible)

	//TODO надо бы подумать как избавиться от этой кольцевой зависимости
	si := NewSiService(&SiDeps{
		Repo:         deps.Repo.SI,
		TxManager:    txManager,
		Instrument:   instrument,
		Verification: verification,
		Location:     NewLocationService(&LocationDeps{Repo: deps.Repo.Location, TxManager: txManager, Responsible: responsible, ActivityLog: activityLog}),
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
	location := NewLocationService(&LocationDeps{
		Repo:         deps.Repo.Location,
		TxManager:    txManager,
		Responsible:  responsible,
		Notification: notification,
		Most:         most,
		ActivityLog:  activityLog,
	})
	department := NewDepartmentService(deps.Repo.Department, location)
	employee := NewEmployeeService(deps.Repo.Employee, location)
	departmentAccess := NewDepartmentAccessService(deps.Repo.DepartmentAccess)
	channel := NewChannelService(deps.Repo.Channel)

	export := NewExportService(&ExportDeps{
		File:    file,
		SI:      si,
		Columns: columns,
	})
	importFile := NewImportService(&ImportDeps{
		Instrument:     instrument,
		Verification:   verification,
		Repair:         repair,
		Preservation:   preservation,
		TransferToSave: transferToSave,
		TransferToDep:  transferToDep,
		WriteOff:       writeOff,
	})

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
		SiStatus:             status,
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
		ActivityLog:          activityLog,
		Filters:              filters,
		Sorting:              sorting,
		Department:           department,
		DepartmentAccess:     departmentAccess,
		Employee:             employee,
		Channel:              channel,
		Responsible:          responsible,
		Location:             location,

		File:         file,
		Notification: notification,
		Scheduler:    scheduler,
		Most:         most,
		Export:       export,
		ImportFile:   importFile,
	}
}
