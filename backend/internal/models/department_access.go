package models

type DepartmentAccess struct {
	Id           string `json:"id" db:"id"`
	DepartmentId string `json:"departmentId" db:"department_id"`
	UserId       string `json:"userId" db:"sso_id"`
}

type GetDepartmentAccessDTO struct {
	DepartmentId string `json:"departmentId"`
	UserId       string `json:"userId"`
	RealmId      string `json:"realmId"`
}

type ReplaceDepartmentAccessDTO struct {
	DepartmentId string   `json:"departmentId" db:"department_id" binding:"required"`
	UserIDs      []string `json:"userIds"`
}

type DepartmentAccessDTO struct {
	Id           string `json:"id" db:"id"`
	DepartmentId string `json:"departmentId" db:"department_id" binding:"required"`
	UserId       string `json:"userId" db:"sso_id" binding:"required"`
}

type DeleteDepartmentAccessDTO struct {
	Id string `json:"id" db:"id"`
}
