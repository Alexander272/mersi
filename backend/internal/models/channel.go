package models

type Channel struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	MostChannelId string `json:"mostChannelId" db:"most_channel_id"`
}

type GetChannel struct{}
