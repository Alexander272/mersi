export interface IResponsible {
	id: string
	departmentId: string
	ssoId: string
}

export interface IChangeResponsible {
	new: IResponsible[]
	updated: IResponsible[]
	deleted: string[]
}
