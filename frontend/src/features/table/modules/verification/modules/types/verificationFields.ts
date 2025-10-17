export interface IVerificationField {
	id: string
	sectionId: string
	field: string
	label: string
	type: string
	position: number
	width: number
}

export interface IVerificationFieldDTO {
	id: string
	sectionId: string
	field: string
	label: string
	type: string
	position: number
	group: string
	width?: number
}

export interface IVerificationFieldForm {
	id: string
	sectionId: string
	field: string
	label: string
	type: string
	position: number
	width: number
	group: string
	status?: 'new' | 'updated' | 'deleted' | 'moved' | 'none'
}
