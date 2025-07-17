package repository

import (
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

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

type Repository struct {
	RuleItem
	Rule
	Role

	Realm
	Accesses
	Section
	Columns
	CreateForm
	Document
	Instrument
	Verification
	VerificationDoc
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
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		RuleItem: postgres.NewRuleItemRepo(db),
		Rule:     postgres.NewRuleRepo(db),
		Role:     postgres.NewRoleRepo(db),

		Realm:                postgres.NewRealmRepo(db),
		Accesses:             postgres.NewAccessesRepo(db),
		Section:              postgres.NewSectionRepo(db),
		Columns:              postgres.NewColumnRepo(db),
		CreateForm:           postgres.NewCreateFormRepo(db),
		Instrument:           postgres.NewInstrumentRepo(db),
		Document:             postgres.NewDocumentRepo(db),
		Verification:         postgres.NewVerificationRepo(db),
		VerificationDoc:      postgres.NewVerificationDocRepo(db),
		SI:                   postgres.NewSIRepo(db),
		ContextMenu:          postgres.NewContextRepo(db),
		CustomContextMenu:    postgres.NewCustomContextRepo(db),
		ToolsMenu:            postgres.NewToolsMenuRepo(db),
		Repair:               postgres.NewRepairRepo(db),
		VerificationFields:   postgres.NewVerificationFieldRepo(db),
		Preservation:         postgres.NewPreservationRepo(db),
		TransferToSave:       postgres.NewTransferToSaveRepo(db),
		TransferToDepartment: postgres.NewTransferToDepRepo(db),
		WriteOff:             postgres.NewWriteOffRepo(db),
	}
}
