export interface ILocation {
	id: string
	instrumentId: string
	department: string
	person: string
	place?: string
	personId?: string
	departmentId?: string
	lastPlaceId?: string
	dateOfIssue: number
	dateOfReceiving: number
	needConfirm: boolean
	hasConfirmed?: boolean
	status: string
}

export interface ILocationDTO {
	id?: string
	instrumentId: string
	person: string
	department: string
	dateOfIssue: number
	dateOfReceiving: number
	needConfirm: boolean
	status: string
	isToReserve?: boolean
}

export interface ILocationForm {
	id?: string
	department: string
	person: string
	dateOfIssue: number
	needConfirm: boolean
	isToReserve?: boolean
}

export interface IReceiving {
	status: string
	instrumentId: string[]
}

export interface IForcedReceipt {
	instrumentId: string
}
