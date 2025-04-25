export interface IColumn {
	id: string
	sectionId: string
	name: string
	field: string
	position: number
	type: string
	width: number
	parentId: string
	allowSort: boolean
	allowFilter: boolean
	created: Date
	children?: IColumn[]
}

export interface IColumnDTO {
	id: string
	sectionId: string
	name: string
	field: string
	position: number
	type: string
	width?: number
	parentId?: string
	allowSort: boolean
	allowFilter: boolean
}

export type ColumnTypes = 'text' | 'number' | 'date' | 'file' | 'list' | 'autocomplete' | 'parent'
