package models

import (
	"mime/multipart"
)

type Document struct {
	Id           string `json:"id" db:"id"`
	Label        string `json:"label" db:"label"`
	Size         int64  `json:"size" db:"size"`
	Path         string `json:"path" db:"path"`
	DocumentType string `json:"type" db:"type"`
	Group        string `json:"group" db:"belongs"`
	InstrumentId string `json:"-" db:"instrument_id"`
	UserId       string `json:"userId" db:"user_id"`
}

type GetDocumentDTO struct {
	UserId       string `json:"userId" db:"user_id"`
	Group        string `json:"group" db:"belongs"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
}

type GetDocumentsDTO struct {
	VerificationId string
	InstrumentId   string
	Group          string `json:"group" db:"belongs"`
}

type DocumentsDTO struct {
	InstrumentId string
	UserId       string
	Group        string
	Files        []*multipart.FileHeader
}

type PathParts struct {
	InstrumentId string
	UserId       string
	Group        string
	IdWasEmpty   bool
}

type DeleteDocumentDTO struct {
	Id           string
	InstrumentId string
	UserId       string
	Group        string
	Filename     string
	IsTemp       bool
}
