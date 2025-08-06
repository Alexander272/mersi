package models

type Department struct {
	Id          string `json:"id" db:"id"`
	RealmId     string `json:"realmId" db:"realm_id"`
	Name        string `json:"name" db:"name"`
	LeaderId    string `json:"leaderId" db:"leader_id"`
	ChannelId   string `json:"channelId" db:"channel_id"`
	ChannelName string `json:"channelName" db:"channel_name"`
}

type GetDepartmentsDTO struct {
	RealmId string `json:"realmId" db:"realm_id"`
}

type GetDepartmentByIdDTO struct {
	Id      string `json:"id" db:"id" binding:"required"`
	RealmId string `json:"realmId" db:"realm_id"`
}

type DepartmentDTO struct {
	Id          string `json:"id" db:"id"`
	RealmId     string `json:"realmId" db:"realm_id"`
	Name        string `json:"name" db:"name" binding:"required"`
	LeaderId    string `json:"leaderId" db:"leader_id"`
	ChannelId   string `json:"channelId" db:"channel_id"`
	ChannelName string `json:"channelName" db:"channel_name"`
}
