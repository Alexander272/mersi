export interface IStatus {
	id: string
	sectionId: string
	position: number
	value: string
	label: string
}

export interface IStatusDTO {
	id: string
	sectionId: string
	position: number
	value: string
	label: string
}

export interface IStatusForm {
	id: string
	position: number
	value: string
	label: string
	status?: 'new' | 'updated' | 'deleted' | 'moved' | 'none'
}
