package models

type Responsible struct {
	Id           string `json:"id" db:"id"`
	DepartmentId string `json:"departmentId" db:"department_id"`
	UserId       string `json:"userId" db:"sso_id"`
}
type ResponsibleWithChannel struct {
	Id           string `json:"id" db:"id"`
	DepartmentId string `json:"departmentId" db:"department_id"`
	UserId       string `json:"userId" db:"sso_id"`
	ChannelId    string `json:"channelId" db:"channel_id"`
}

type ResponsibleDTO struct {
	Id           string `json:"id" db:"id"`
	DepartmentId string `json:"departmentId" db:"department_id"`
	UserId       string `json:"userId" db:"sso_id"`
}

type GetResponsibleDTO struct {
	DepartmentId string `json:"departmentId" db:"department_id"`
	UserId       string `json:"userId" db:"sso_id"`
}

type ChangeResponsibleDTO struct {
	New     []*ResponsibleDTO `json:"new"`
	Updated []*ResponsibleDTO `json:"updated"`
	Deleted []string          `json:"deleted"`
}
