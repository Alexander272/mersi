package models

import "time"

type SI struct {
	Id                        string           `json:"id" db:"id"`
	Position                  int              `json:"position" db:"position"`
	Name                      string           `json:"name" db:"name"`
	DateOfReceipt             time.Time        `json:"dateOfReceipt" db:"date_of_receipt"`
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
	VerificationDate          time.Time        `json:"verificationDate" db:"date"`
	NextVerificationDate      time.Time        `json:"nextVerificationDate" db:"next_date"`
	Certificate               string           `json:"certificate" db:"certificate"`
	CertificateId             string           `json:"certificateId" db:"certificate_id"`
	LastCertificate           string           `json:"lastCertificate" db:"last_certificate"`
	LastCertificateId         string           `json:"lastCertificateId" db:"last_certificate_id"`
	RepairInfo                string           `json:"repairInfo" db:"repair"`
	PreservationDate          time.Time        `json:"preservationDate" db:"preservation"`
	DePreservationDate        time.Time        `json:"dePreservationDate" db:"de_preservation"`
	TransferDate              time.Time        `json:"transferDate" db:"transfer_date"`
	ReturnDate                time.Time        `json:"returnDate" db:"return_date"`
	TransferToDepartment      string           `json:"transferToDepartment" db:"transfer_to_dep"`
	WriteOff                  string           `json:"writeOff" db:"write_off"`
	Person                    string           `json:"person" db:"person"`
	Place                     string           `json:"place" db:"place"`
	LastPlace                 string           `json:"lastPlace" db:"last_place"`
	Status                    InstrumentStatus `json:"status" db:"status"`
	Total                     int              `json:"total" db:"total"`
}

type SiWithLog struct {
	Id               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	DateOfReceipt    time.Time `json:"dateOfReceipt" db:"date_of_receipt"`
	Type             string    `json:"type" db:"type"`
	FactoryNumber    string    `json:"factoryNumber" db:"factory_number"`
	Responsible      string    `json:"responsible" db:"responsible"`
	RepairInfo       string    `json:"repairInfo" db:"repair"`
	PreservationInfo string    `json:"preservationInfo" db:"preservation"`
	SavingInfo       string    `json:"savingInfo" db:"saving"`
	WriteOff         string    `json:"writeOff" db:"write_off"`
}

type BaseSI struct {
	Instrument   *Instrument   `json:"instrument"`
	Verification *Verification `json:"verification"`
	Location     *Location     `json:"location"`
}

type GetSiDTO struct {
	SectionId        string
	Page             *Page
	Sort             []*Sort
	Filters          []*Filter
	Search           *Search
	Status           InstrumentStatus
	UserID           string
	DepartmentAccess []string
}

type GetSiByIdDTO struct {
	Id string
}

type SiDTO struct {
	Instrument   *InstrumentDTO   `json:"instrument" binding:"required"`
	Verification *VerificationDTO `json:"verification"`
	Location     *LocationDTO     `json:"location"`
}

type ChangePositionDTO struct {
	SectionId   string `json:"sectionId" db:"section_id" binding:"required"`
	NewPosition int    `json:"newPosition" db:"new_position" binding:"min=0"`
	OldPosition int    `json:"oldPosition" db:"old_position" binding:"min=0"`
}

type SiVerification struct {
	SI                  []*SI
	NotificationChannel string `json:"notification" db:"notification_channel"`
	BidType             string `json:"bidType" db:"bid_type"`
}

type SiReceiving struct {
	PostId  string `json:"-"`
	Status  string `json:"status"`
	Channel string `json:"channel"`
	Place   string `json:"place"`
	SI      []*SI  `json:"si"`
}

type DeleteSiDTO struct {
	Id    string
	Actor *Actor
}
