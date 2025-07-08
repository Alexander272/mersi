package models

import "time"

type Preservation struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	DateStart    int64     `json:"dateStart" db:"date_start"`
	DateEnd      int64     `json:"dateEnd" db:"date_end"`
	NotesStart   string    `json:"notesStart" db:"notes_start"`
	NotesEnd     string    `json:"notesEnd" db:"notes_end"`
	Created      time.Time `json:"created" db:"created_at"`
}

type GetPreservationsDTO struct {
	InstrumentId string `json:"instrumentId"`
}

type PreservationDTO struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
	DateStart    int64  `json:"dateStart" db:"date_start"`
	DateEnd      int64  `json:"dateEnd" db:"date_end"`
	NotesStart   string `json:"notesStart" db:"notes_start"`
	NotesEnd     string `json:"notesEnd" db:"notes_end"`
	Notes        string `json:"notes" db:"notes"`
}

type DeletePreservationDTO struct {
	Id string `json:"id" db:"id"`
}
