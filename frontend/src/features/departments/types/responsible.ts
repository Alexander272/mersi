export interface IResponsible {
	id: string
	departmentId: string
	userId: string
}

export interface IChangeResponsible {
	new: IResponsible[]
	updated: IResponsible[]
	deleted: string[]
}
