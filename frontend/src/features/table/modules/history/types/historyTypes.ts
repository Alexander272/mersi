export interface IHistoryType {
	id: string
	sectionId: string
	group: string
	label: string
	position: number
	created: string
}

export interface IHistoryTypeDTO {
	id: string
	sectionId: string
	group: string
	label: string
	position: number
}

export interface IHistoryForm {
	id: string
	sectionId: string
	group: string
	label: string
	position: number
	status?: 'new' | 'updated' | 'deleted' | 'moved' | 'none'
}
