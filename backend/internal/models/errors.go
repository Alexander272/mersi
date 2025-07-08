package models

import "errors"

var (
	ErrNoRows        = errors.New("row not found")
	ErrNoChannel     = errors.New("channel not found")
	ErrNoResponsible = errors.New("responsible not found")
	ErrNoInstrument  = errors.New("instrument not found")

	ErrSessionEmpty = errors.New("user session not found")

	ErrNotValid = errors.New("data is not valid")
)
