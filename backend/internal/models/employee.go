package models

type GetEmployeesDTO struct {
	Filters map[string]string
	// DepartmentId string `json:"departmentId"`
}

// type GetDepartmentsDTO struct {
// 	UserId string `json:"userId" binding:"required"`
// }

type Employee struct {
	Id           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Notes        string `json:"notes" db:"notes"`
	DepartmentId string `json:"departmentId" db:"department_id"`
	// IncludeDepartments []string `json:"includeDepartments" db:"include_departments"` //TODO возможно это нужно делать в другом месте
	// DefaultFilters     string   // TODO определиться с типом и форматом фильтров (возможно фильтры нужно делать в другом месте)
	MattermostId string `json:"mattermostId" db:"most_id"` //TODO возможно на клиенте это мне не нужно будет
	// Type string => member or leader
	IsLead bool `json:"isLead" db:"is_lead"`
	//TODO нужно как-то связывать работника и пользователя (чтобы пользователь мог вернуть инструмент) можно сгенерировать какой-нибудь код, записать его в keycloak и бд, а потом по нему искать запись
}

type EmployeeData struct {
	Id         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Notes      string `json:"notes" db:"notes"`
	Department string `json:"department" db:"department"`
	MostId     string `json:"mostId" db:"most_id"`
	IsLead     bool   `json:"isLead" db:"is_lead"`
}

type WriteEmployeeDTO struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Notes        string `json:"notes" db:"notes"`
	DepartmentId string `json:"departmentId"`
	MattermostId string `json:"mattermostId"`
}
