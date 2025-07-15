package models

type SI struct {
	Id                        string `json:"id" db:"id"`
	Position                  int    `json:"position" db:"position"`
	Name                      string `json:"name" db:"name"`
	DateOfReceipt             int    `json:"dateOfReceipt" db:"date_of_receipt"`
	Type                      string `json:"type" db:"type"`
	FactoryNumber             string `json:"factoryNumber" db:"factory_number"`
	MeasurementLimits         string `json:"measurementLimits" db:"measurement_limits"`
	Accuracy                  string `json:"accuracy" db:"accuracy"`
	StateRegister             string `json:"stateRegister" db:"state_register"`
	CountryOfProduce          string `json:"countryOfProduce" db:"country_of_produce"`
	Manufacturer              string `json:"manufacturer" db:"manufacturer"`
	Responsible               string `json:"responsible" db:"responsible"`
	Inventory                 string `json:"inventory" db:"inventory"`
	YearOfIssue               int    `json:"yearOfIssue" db:"year_of_issue"`
	InterVerificationInterval int    `json:"interVerificationInterval" db:"inter_verification_interval"`
	ActOfEntering             string `json:"actOfEntering" db:"act_of_entering"`
	ActOfEnteringId           string `json:"actOfEnteringId" db:"act_of_entering_id"`
	Notes                     string `json:"notes" db:"notes"`
	VerificationDate          int    `json:"verificationDate" db:"date"`
	NextVerificationDate      int    `json:"nextVerificationDate" db:"next_date"`
	Certificate               string `json:"certificate" db:"certificate"`
	CertificateId             string `json:"certificateId" db:"certificate_id"`
	Repair                    string `json:"repairInfo" db:"repair"`
	Preservation              int    `json:"preservationDate" db:"preservation"`
	DePreservation            int    `json:"dePreservationDate" db:"de_preservation"`
	TransferDate              int    `json:"transferDate" db:"transfer_date"`
	ReturnDate                int    `json:"returnDate" db:"return_date"`
	//TODO дописать оставшиеся поля

	Total int `json:"total" db:"total"`
}

type BaseSI struct {
	Instrument   *Instrument   `json:"instrument"`
	Verification *Verification `json:"verification"`
}

type GetSiDTO struct {
	SectionId string
	Page      *Page
	Sort      []*Sort
	Filters   []*Filter
	Search    *Search
	Status    InstrumentStatus
}

type GetSiByIdDTO struct {
	Id string
}

type SiDTO struct {
	Instrument   *InstrumentDTO   `json:"instrument" binding:"required"`
	Verification *VerificationDTO `json:"verification"`
}

type ChangePositionDTO struct {
	SectionId   string `json:"sectionId" db:"section_id" binding:"required"`
	NewPosition int    `json:"newPosition" db:"new_position" binding:"min=0"`
	OldPosition int    `json:"oldPosition" db:"old_position" binding:"min=0"`
}
