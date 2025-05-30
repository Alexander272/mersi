package models

type SiDTO struct {
	Instrument   *InstrumentDTO   `json:"instrument" binding:"required"`
	Verification *VerificationDTO `json:"verification"`
}
