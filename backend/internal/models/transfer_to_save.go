package models

import "time"

type TransferToSave struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	DateStart    time.Time `json:"dateStart" db:"date_start"`
	DateEnd      time.Time `json:"dateEnd" db:"date_end"`
	NotesStart   string    `json:"notesStart" db:"notes_start"`
	NotesEnd     string    `json:"notesEnd" db:"notes_end"`
	Created      time.Time `json:"created" db:"created_at"`
}

type GetTransferToSaveDTO struct {
	InstrumentId string `json:"instrumentId"`
}

type TransferToSaveDTO struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	DateStart    time.Time `json:"dateStart" db:"date_start"`
	DateEnd      time.Time `json:"dateEnd" db:"date_end"`
	NotesStart   string    `json:"notesStart" db:"notes_start"`
	NotesEnd     string    `json:"notesEnd" db:"notes_end"`
}

type DeleteTransferToSaveDTO struct {
	Id string `json:"id" db:"id"`
}
