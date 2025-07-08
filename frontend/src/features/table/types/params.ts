import { ColumnTypes } from '@/features/sections/modules/columns/types/columns'

export interface IParams {
	page?: number
	size?: number
	sort?: ISort
	filters?: IFilter[]
}

export type ISort = {
	[x: string]: 'DESC' | 'ASC'
}

export type CompareTypes = 'con' | 'start' | 'end' | 'like' | 'in' | 'eq' | 'gte' | 'lte' | 'range' | 'null'
export interface IFilter {
	field: string
	fieldType: Exclude<ColumnTypes, 'parent' | 'file'>
	compareType: CompareTypes
	value: string
}

export interface ISearch {
	value: string
	fields: string[]
}
