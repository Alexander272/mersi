export interface IDepartment {
	id: string
	name: string
	leaderId: string
	channelId: string
	channelName: string
}

export interface IDepartmentAccess {
	id: string
	departmentId: string
	userId: string
}

export interface IChangeDepAccess {
	departmentId: string
	userIds: string[]
}
