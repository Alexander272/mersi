package repository

import (
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Transaction interface {
	postgres.Transaction
}

type RuleItem interface {
	postgres.RuleItem
}
type Rule interface {
	postgres.Rule
}
type Role interface {
	postgres.Role
}

type Realm interface {
	postgres.Realm
}
type Accesses interface {
	postgres.Accesses
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
type Instrument interface {
	postgres.Instrument
}
type Document interface {
	postgres.Document
}
type Verification interface {
	postgres.Verification
}
type VerificationDoc interface {
	postgres.VerificationDoc
}
type Location interface {
	postgres.Location
}
type SI interface {
	postgres.SI
}
type ContextMenu interface {
	postgres.ContextMenu
}
type CustomContextMenu interface {
	postgres.CustomContextMenu
}
type ToolsMenu interface {
	postgres.ToolsMenu
}
type SiStatus interface {
	postgres.Status
}
type Repair interface {
	postgres.Repair
}
type VerificationFields interface {
	postgres.VerificationFields
}
type Preservation interface {
	postgres.Preservation
}
type TransferToSave interface {
	postgres.TransferToSave
}
type TransferToDepartment interface {
	postgres.TransferToDepartment
}
type WriteOff interface {
	postgres.WriteOff
}
type HistoryType interface {
	postgres.HistoryType
}
type Filters interface {
	postgres.Filters
}
type Sorting interface {
	postgres.Sorting
}
type Department interface {
	postgres.Department
}
type DepartmentAccess interface {
	postgres.DepartmentAccess
}
type Employee interface {
	postgres.Employee
}
type Channel interface {
	postgres.Channel
}
type Responsible interface {
	postgres.Responsible
}
type Users interface {
	postgres.User
}

type Repository struct {
	Transaction

	RuleItem
	Rule
	Role

	Realm
	Accesses
	Section
	Columns
	CreateForm
	SiStatus
	Document
	Instrument
	Verification
	VerificationDoc
	Location
	SI
	ContextMenu
	CustomContextMenu
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
	Department
	DepartmentAccess
	Employee
	Channel
	Responsible
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	transactions := postgres.NewTransactionRepo(db)

	return &Repository{
		Transaction: transactions,

		RuleItem: postgres.NewRuleItemRepo(db),
		Rule:     postgres.NewRuleRepo(db),
		Role:     postgres.NewRoleRepo(db),

		Realm:                postgres.NewRealmRepo(db),
		Accesses:             postgres.NewAccessesRepo(db),
		Section:              postgres.NewSectionRepo(db),
		Columns:              postgres.NewColumnRepo(db),
		CreateForm:           postgres.NewCreateFormRepo(db),
		SiStatus:             postgres.NewStatusRepo(db),
		Instrument:           postgres.NewInstrumentRepo(db, transactions),
		Document:             postgres.NewDocumentRepo(db),
		Verification:         postgres.NewVerificationRepo(db, transactions),
		VerificationDoc:      postgres.NewVerificationDocRepo(db, transactions),
		Location:             postgres.NewLocationRepo(db, transactions),
		SI:                   postgres.NewSIRepo(db, transactions),
		ContextMenu:          postgres.NewContextRepo(db),
		CustomContextMenu:    postgres.NewCustomContextRepo(db),
		ToolsMenu:            postgres.NewToolsMenuRepo(db),
		Repair:               postgres.NewRepairRepo(db, transactions),
		VerificationFields:   postgres.NewVerificationFieldRepo(db),
		Preservation:         postgres.NewPreservationRepo(db, transactions),
		TransferToSave:       postgres.NewTransferToSaveRepo(db, transactions),
		TransferToDepartment: postgres.NewTransferToDepRepo(db, transactions),
		WriteOff:             postgres.NewWriteOffRepo(db, transactions),
		HistoryType:          postgres.NewHistoryTypeRepo(db),
		Filters:              postgres.NewFilterRepo(db),
		Sorting:              postgres.NewSortingRepo(db),
		Department:           postgres.NewDepartmentRepo(db),
		DepartmentAccess:     postgres.NewDepartmentAccessRepo(db),
		Employee:             postgres.NewEmployeeRepo(db),
		Channel:              postgres.NewChannelRepo(db),
		Responsible:          postgres.NewResponsibleRepo(db),
		Users:                postgres.NewUserRepo(db, transactions),
	}
}
