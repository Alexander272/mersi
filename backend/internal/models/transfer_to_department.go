package models

import "time"

type TransferToDepartment struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	Date         int64     `json:"date" db:"date"`
	Notes        string    `json:"notes" db:"notes"`
	DocId        string    `json:"docId" db:"doc_id"`
	DocName      string    `json:"docName" db:"doc_name"`
	Created      time.Time `json:"created" db:"created_at"`
}

type GetTransferToDepDTO struct {
	InstrumentId string `json:"instrumentId"`
}

type TransferToDepartmentDTO struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
	Date         int64  `json:"date" db:"date"`
	Notes        string `json:"notes" db:"notes"`
	DocId        string `json:"docId" db:"doc_id"`
	DocName      string `json:"docName" db:"doc_name"`
}

type DeleteTransferToDepDTO struct {
	Id string `json:"id" db:"id"`
}
