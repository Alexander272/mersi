package models

import "time"

type Location struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
	Person       string `json:"person" db:"person"`
	// Department      string `json:"department" db:"department"`
	Place           string    `json:"place" db:"place"`
	PersonId        string    `json:"personId" db:"person_id"`
	DepartmentId    string    `json:"departmentId" db:"department_id"`
	DateOfIssue     time.Time `json:"dateOfIssue" db:"date_of_issue"`
	DateOfReceiving time.Time `json:"dateOfReceiving" db:"date_of_receiving"`
	NeedConfirm     bool      `json:"needConfirm" db:"need_confirmed"`
	HasConfirmed    bool      `json:"hasConfirmed" db:"has_confirmed"`
	LastPlace       string    `json:"lastPlace" db:"last_place"`
	LastPlaceId     string    `json:"lastPlaceId" db:"last_place_id"`
	Status          string    `json:"status" db:"status"`
}

type GetLocationDTO struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
}

type GetSeveralLocationsDTO struct {
	InstrumentIds []string `json:"instrumentIds" db:"instrument_id"`
}

type GetLocationByHolderDTO struct {
	PersonId string `json:"personId" db:"person_id"`
}
type GetLocationByDepartmentDTO struct {
	DepartmentId string `json:"departmentId" db:"department_id"`
}

type LocationDTO struct {
	Id           string `json:"id" db:"id"`
	Actor        *Actor
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
	PersonId     string `json:"person" db:"person_id"`
	DepartmentId string `json:"department" db:"department_id"`
	// ReceiptDate     string `json:"receiptDate" db:"receipt_date"`
	// DeliveryDate    string `json:"deliveryDate" db:"delivery_date"`
	DateOfIssue     time.Time `json:"dateOfIssue" db:"date_of_issue" binding:"required"`
	DateOfReceiving time.Time `json:"dateOfReceiving" db:"date_of_receiving"`
	NeedConfirm     bool      `json:"needConfirm" db:"need_confirmed"`
	Status          string    `json:"status" db:"status"`
	UserId          string    `json:"userId" db:"user_id"`
}

type DeleteLocationDTO struct {
	Id    string `json:"id" db:"id"`
	Actor *Actor
}

type ReceivingDTO struct {
	InstrumentIds []string `json:"instrumentId" db:"instrument_id"`
	Status        string   `json:"status" db:"status"` // либо отправляется в резерв, либо к сотруднику
	UserId        string
	HasConfirmed  bool
	// Missing       []SelectedSI
	// DateOfReceiving string   `json:"dateOfReceiving" db:"date_of_receiving"`
}

type ForcedReceiptDTO struct {
	InstrumentId string `json:"instrumentId" binding:"required"`
}

type SelectByDepsDTO struct {
	DepartmentIds []string
	InstrumentIds []string
	Status        string
}
