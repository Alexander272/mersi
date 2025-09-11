export interface ILocation {
	id: string
	instrumentId: string
	department: string
	person: string
	place?: string
	personId?: string
	departmentId?: string
	lastPlaceId?: string
	dateOfIssue: string
	dateOfReceiving: string
	needConfirm: boolean
	hasConfirmed?: boolean
	status: string
}

export interface ILocationDTO {
	id?: string
	instrumentId: string
	person: string
	department: string
	dateOfIssue: string
	dateOfReceiving: string
	needConfirm: boolean
	status: string
	isToReserve?: boolean
}

export interface ILocationForm {
	id?: string
	department: string
	person: string
	dateOfIssue: string
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
