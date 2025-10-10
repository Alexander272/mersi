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
