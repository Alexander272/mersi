package models

import "time"

type WriteOff struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	Date         time.Time `json:"date" db:"date"`
	Notes        string    `json:"notes" db:"notes"`
	DocId        string    `json:"docId" db:"doc_id"`
	DocName      string    `json:"docName" db:"doc_name"`
	Created      time.Time `json:"created" db:"created_at"`
}

type GetWriteOffDTO struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId"`
}

type WriteOffDTO struct {
	Id           string `json:"id" db:"id"`
	Actor        *Actor
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	Date         time.Time `json:"date" db:"date"`
	Notes        string    `json:"notes" db:"notes"`
	DocId        string    `json:"docId" db:"doc_id"`
	DocName      string    `json:"docName" db:"doc_name"`
	UserId       string    `json:"userId" db:"user_id"`
	DeletedDocs  []DeletedDoc `json:"deletedDocs"`
}

type DeleteWriteOffDTO struct {
	Id    string `json:"id" db:"id"`
	Actor *Actor
}
