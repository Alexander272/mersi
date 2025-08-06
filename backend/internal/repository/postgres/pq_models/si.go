package pq_models

import "github.com/Alexander272/mersi/backend/internal/models"

type SI struct {
	models.SI
	NotificationChannel string `json:"notification" db:"notification_channel"`
	BidType             string `json:"bidType" db:"bid_type"`
}
