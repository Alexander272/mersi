package models

import "time"

type Repair struct {
	Id           string    `json:"id" db:"id"`
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	Defect       string    `json:"defect" db:"defect"`
	Work         string    `json:"work" db:"work"`
	PeriodStart  time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd    time.Time `json:"periodEnd" db:"period_end"`
	Description  string    `json:"description" db:"description"`
	Created      time.Time `json:"created" db:"created_at"`
}

type GetRepairDTO struct {
	InstrumentId string `json:"instrumentId"`
}

type RepairDTO struct {
	Id           string `json:"id" db:"id"`
	Actor        *Actor
	InstrumentId string    `json:"instrumentId" db:"instrument_id"`
	Defect       string    `json:"defect" db:"defect"`
	Work         string    `json:"work" db:"work"`
	PeriodStart  time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd    time.Time `json:"periodEnd" db:"period_end"`
	Description  string    `json:"description" db:"description"`
}

type DeleteRepairDTO struct {
	Id    string `json:"id"`
	Actor *Actor
}
