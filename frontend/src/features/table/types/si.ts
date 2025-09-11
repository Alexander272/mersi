import type { ILocationDTO } from '../modules/locations/types/location'
import type { IVerificationDTO } from '../modules/verification/types/verification'
import type { IInstrumentDTO } from './instrument'
import type { ISort, IFilter, ISearch } from './params'

export interface ISI {
	id: string
	position: number
	name: string
	dateOfReceipt: string
	type: string
	factoryNumber: string
	measurementLimits: string
	accuracy: string
	stateRegister: string
	countryOfProduce: string
	manufacturer: string
	responsible: string
	inventory: string
	yearOfIssue: number
	interVerificationInterval: number
	actOfEntering: string
	actOfEnteringId: string
	notes: string
	verificationDate: string
	nextVerificationDate: string
	certificate: string
	certificateId: string
	repair: string
	preservationDate: string
	dePreservationDate: string
	transferDate: string
	returnDate: string
	transferToDepartment: string
	writeOff: string
	person: string
	place: string
	lastPlace: string
	status: LocationStatus
}
export type LocationStatus = 'reserve' | 'used' | 'moved'

export interface ISiForm {
	instrument: IInstrumentDTO
	verification: IVerificationDTO
	location: ILocationDTO
}

export interface ISiDTO {
	instrument: IInstrumentDTO
	verification?: IVerificationDTO
}

export type Status = 'work' | 'repair' | 'decommissioning' | 'transferred'
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

export interface ISentSI {
	place: string
	si: ISI[]
	channel: string
	status: string
}
