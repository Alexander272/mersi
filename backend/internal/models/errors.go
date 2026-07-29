package models

import "errors"

var (
	ErrNoRows   = errors.New("row not found")
	ErrBadValue = errors.New("bad value")

	ErrSessionEmpty = errors.New("user session not found")

	ErrInvalidInput      = errors.New("invalid input data")
	ErrInvalidPermission = errors.New("invalid permission")

	ErrReservedRole          = errors.New("cannot create or update reserved role")
	ErrCircularInheritance   = errors.New("circular inheritance detected")
	ErrCannotInheritFromSelf = errors.New("role cannot inherit from itself")
	ErrParentRoleNotFound    = errors.New("parent role not found or inactive")
	ErrRoleNotEditable       = errors.New("role is not editable")

	ErrHasInstrument    = errors.New("instrument is still in use")
	ErrAlreadyExists    = errors.New("record already exists")
	ErrNoResponsible    = errors.New("no responsible person")
	ErrNoChannel        = errors.New("no channel found")
	ErrNoInstrument     = errors.New("instrument not found")
	ErrNotValid         = errors.New("data is not valid")
)
