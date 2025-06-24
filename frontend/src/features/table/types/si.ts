import { IInstrumentDTO } from './instrument'
import { ISort, IFilter } from './params'
import { IVerificationDTO } from './verification'

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
	filters?: IFilter[]
}
