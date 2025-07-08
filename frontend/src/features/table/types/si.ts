import type { IVerificationDTO } from '../modules/verification/types/verification'
import type { IInstrumentDTO } from './instrument'
import type { ISort, IFilter, ISearch } from './params'

export interface ISI {
	id: string
}

export interface ISiForm {
	instrument: IInstrumentDTO
	verification: IVerificationDTO
}

export interface ISiDTO {
	instrument: IInstrumentDTO
	verification?: IVerificationDTO
}

export type Status = 'work' | 'repair' | 'decommissioning'
export interface IGetSiDTO {
	section: string
	status?: Status
	all?: boolean
	page?: number
	size?: number
	sort?: ISort
	search?: ISearch
	filters?: IFilter[]
}

export interface IChangePositionDTO {
	sectionId: string
	newPosition: number
	oldPosition: number
}
