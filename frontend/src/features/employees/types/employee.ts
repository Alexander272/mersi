export interface IDepartment {
	id: string
	name: string
	leaderId: string
}

export interface IEmployee {
	id: string
	name: string
	departmentId: string
	isLead?: boolean
	mattermostId?: string
	notes: string
}

export interface IEmployeeForm {
	lastName: string
	firstName: string
	surname: string
	departmentId: string
	isLead: boolean
}

export interface IDepartmentForm {
	name: string
}
