package models

type InstrumentStatus string

const (
	InstrumentStatusWork   InstrumentStatus = "work"
	InstrumentStatusRepair InstrumentStatus = "repair"
	InstrumentStatusDec    InstrumentStatus = "decommissioning"
	InstrumentDeleted      InstrumentStatus = "deleted"
	InstrumentDraft        InstrumentStatus = "draft"
)

type Instrument struct {
	Id                        string           `json:"id" db:"id"`
	Position                  int              `json:"position" db:"position"`
	Name                      string           `json:"name" db:"name"`
	DateOfReceipt             int64            `json:"dateOfReceipt" db:"date_of_receipt"`
	Type                      string           `json:"type" db:"type"`
	FactoryNumber             string           `json:"factoryNumber" db:"factory_number"`
	MeasurementLimits         string           `json:"measurementLimits" db:"measurement_limits"`
	Accuracy                  string           `json:"accuracy" db:"accuracy"`
	StateRegister             string           `json:"stateRegister" db:"state_register"`
	CountryOfProduce          string           `json:"countryOfProduce" db:"country_of_produce"`
	Manufacturer              string           `json:"manufacturer" db:"manufacturer"`
	Responsible               string           `json:"responsible" db:"responsible"`
	Inventory                 string           `json:"inventory" db:"inventory"`
	YearOfIssue               int              `json:"yearOfIssue" db:"year_of_issue"`
	InterVerificationInterval int              `json:"interVerificationInterval" db:"inter_verification_interval"`
	ActOfEntering             string           `json:"actOfEntering" db:"act_of_entering"`
	ActOfEnteringId           string           `json:"actOfEnteringId" db:"act_of_entering_id"`
	Notes                     string           `json:"notes" db:"notes"`
	Status                    InstrumentStatus `json:"status" db:"status"`
}

type GetInstrumentsDTO struct {
	SectionId string `json:"sectionId" db:"section_id"`
}

type GetInstrumentByIdDTO struct {
	Id string `json:"id" binding:"required"`
}

type GetUniqueDTO struct {
	Field string `json:"field"`
}

type InstrumentDTO struct {
	Id                        string           `json:"id" db:"id"`
	SectionId                 string           `json:"sectionId" db:"section_id"`
	UserId                    string           `json:"userId" db:"user_id"`
	Position                  int              `json:"position" db:"position" binding:"required"`
	Name                      string           `json:"name" db:"name" binding:"required"`
	DateOfReceipt             int64            `json:"dateOfReceipt" db:"date_of_receipt"`
	Type                      string           `json:"type" db:"type"`
	FactoryNumber             string           `json:"factoryNumber" db:"factory_number"`
	MeasurementLimits         string           `json:"measurementLimits" db:"measurement_limits"`
	Accuracy                  string           `json:"accuracy" db:"accuracy"`
	StateRegister             string           `json:"stateRegister" db:"state_register"`
	CountryOfProduce          string           `json:"countryOfProduce" db:"country_of_produce"`
	Manufacturer              string           `json:"manufacturer" db:"manufacturer"`
	Responsible               string           `json:"responsible" db:"responsible"`
	Inventory                 string           `json:"inventory" db:"inventory"`
	YearOfIssue               int              `json:"yearOfIssue" db:"year_of_issue" binding:"required"`
	InterVerificationInterval int              `json:"interVerificationInterval" db:"inter_verification_interval"`
	ActOfEntering             string           `json:"actOfEntering" db:"act_of_entering"`
	ActOfEnteringId           string           `json:"actOfEnteringId" db:"act_of_entering_id"`
	Notes                     string           `json:"notes" db:"notes"`
	Status                    InstrumentStatus `json:"status" db:"status"`
}

type UpdateStatus struct {
	Id     string           `json:"id" binding:"required"`
	Status InstrumentStatus `json:"status"`
}
