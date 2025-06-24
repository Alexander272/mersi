export interface IColumn {
	id: string
	sectionId: string
	name: string
	field: string
	position: number
	type: ColumnTypes
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
	type: ColumnTypes
	width?: number
	parentId?: string
	allowSort: boolean
	allowFilter: boolean
}

export interface IColumnPositionDTO {
	id: string
	position: number
	parentId: string
}

export type ColumnTypes = 'text' | 'number' | 'date' | 'file' | 'list' | 'autocomplete' | 'parent'
