package pq_models

import (
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type SI struct {
	models.SI
	RepairStart         time.Time `json:"repairStart" db:"repair_start"`
	RepairEnd           time.Time `json:"repairEnd" db:"repair_end"`
	RepairWork          string    `json:"repairWork" db:"repair_work"`
	NotificationChannel string    `json:"notification" db:"notification_channel"`
	BidType             string    `json:"bidType" db:"bid_type"`
}

func (s *SI) ToModel() *models.SI {
	return &models.SI{
		Id:                        s.Id,
		Position:                  s.Position,
		Name:                      s.Name,
		DateOfReceipt:             s.DateOfReceipt,
		Type:                      s.Type,
		FactoryNumber:             s.FactoryNumber,
		MeasurementLimits:         s.MeasurementLimits,
		Accuracy:                  s.Accuracy,
		StateRegister:             s.StateRegister,
		CountryOfProduce:          s.CountryOfProduce,
		Manufacturer:              s.Manufacturer,
		Responsible:               s.Responsible,
		Inventory:                 s.Inventory,
		YearOfIssue:               s.YearOfIssue,
		InterVerificationInterval: s.InterVerificationInterval,
		ActOfEntering:             s.ActOfEntering,
		ActOfEnteringId:           s.ActOfEnteringId,
		Notes:                     s.Notes,
		VerificationDate:          s.VerificationDate,
		NextVerificationDate:      s.NextVerificationDate,
		Certificate:               s.Certificate,
		CertificateId:             s.CertificateId,
		PreservationDate:          s.PreservationDate,
		DePreservationDate:        s.DePreservationDate,
		TransferDate:              s.TransferDate,
		ReturnDate:                s.ReturnDate,
		TransferToDepartment:      s.TransferToDepartment,
		WriteOff:                  s.WriteOff,
		Person:                    s.Person,
		Place:                     s.Place,
		LastPlace:                 s.LastPlace,
		Status:                    s.Status,
		Total:                     s.Total,
	}
}
