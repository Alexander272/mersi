package models

type Verification struct {
	Id           string             `json:"id" db:"id"`
	InstrumentId string             `json:"instrumentId" db:"instrument_id"`
	Date         int64              `json:"verificationDate" db:"date"`
	NextDate     int64              `json:"nextVerificationDate" db:"next_date"`
	RegisterLink string             `json:"registerLink" db:"register_link"`
	NotVerified  bool               `json:"notVerified" db:"not_verified"`
	Notes        string             `json:"notes" db:"notes"`
	Status       string             `json:"status" db:"status"`
	Docs         []*VerificationDoc `json:"docs"`
}

type GetVerificationDTO struct {
	InstrumentId string `json:"instrumentId" db:"instrument_id"`
}

type VerificationDTO struct {
	Id           string                `json:"id" db:"id"`
	InstrumentId string                `json:"instrumentId" db:"instrument_id"`
	Date         int64                 `json:"verificationDate" db:"date"`
	NextDate     int64                 `json:"nextVerificationDate" db:"next_date"`
	RegisterLink string                `json:"registerLink" db:"register_link"`
	NotVerified  bool                  `json:"notVerified" db:"not_verified"`
	Notes        string                `json:"notes" db:"notes"`
	Status       string                `json:"status" db:"status"`
	Docs         []*VerificationDocDTO `json:"docs"`
}

type DeleteVerificationDTO struct {
	Id string `json:"id" db:"id"`
}
