export interface ICreateFormStep {
	step: number
	stepName: string
	fields: ICreateFormField[]
}

export interface ICreateFormField {
	id: string
	sectionId: string
	step: number
	stepName: string
	field: string
	fieldName: string
	path: string
	type: string
	isRequired: boolean
	position: number
	created: string
}

export interface ICreateFormFieldDTO {
	id: string
	sectionId: string
	step: number
	stepName: string
	field: string
	fieldName: string
	path: string
	type: string
	isRequired: boolean
	position: number
}
