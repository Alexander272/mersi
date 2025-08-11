package models

type Location struct {
	Id              string `json:"id" db:"id"`
	InstrumentId    string `json:"instrumentId" db:"instrument_id"`
	Person          string `json:"person" db:"person"`
	Department      string `json:"department" db:"department"`
	Place           string `json:"place" db:"place"`
	PersonId        string `json:"personId" db:"person_id"`
	DepartmentId    string `json:"departmentId" db:"department_id"`
	DateOfIssue     int64  `json:"dateOfIssue" db:"date_of_issue"`
	DateOfReceiving int64  `json:"dateOfReceiving" db:"date_of_receiving"`
	NeedConfirmed   bool   `json:"needConfirmed" db:"need_confirmed"`
	HasConfirmed    bool   `json:"hasConfirmed" db:"has_confirmed"`
	LastPlace       string `json:"lastPlace" db:"last_place"`
	Status          string `json:"status" db:"status"`
}

type GetLocationDTO struct {
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
}

type LocationDTO struct {
	Id           string `json:"id" db:"id"`
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
	PersonId     string `json:"person" db:"person_id"`
	DepartmentId string `json:"department" db:"department_id"`
	// ReceiptDate     string `json:"receiptDate" db:"receipt_date"`
	// DeliveryDate    string `json:"deliveryDate" db:"delivery_date"`
	DateOfIssue     int64  `json:"dateOfIssue" db:"date_of_issue" binding:"required,gte=1000000"`
	DateOfReceiving int64  `json:"dateOfReceiving" db:"date_of_receiving" binding:"gte=0"`
	NeedConfirmed   bool   `json:"needConfirmed" db:"need_confirmed"`
	Status          string `json:"status" db:"status"`
	UserId          string `json:"userId" db:"user_id"`
}

type DeleteLocationDTO struct {
	Id string `json:"id" db:"id"`
}
