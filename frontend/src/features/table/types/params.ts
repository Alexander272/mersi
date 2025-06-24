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
	fieldType: string
	compareType: CompareTypes
	value: string
}
