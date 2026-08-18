type Locations = 'department' | 'place'

export interface IRealm {
	id: string
	name: string
	realm: string
	isActive: boolean
	notificationChannel: string
	expirationNotice: boolean
	returnNotice: boolean
	locationType: Locations
	needConfirmed: boolean
	hasResponsible: boolean
	hasEmployees: boolean
	hasCommissioningCert: boolean
	hasPreservations: boolean
	hasTransfer: boolean
	verificationSubtractDay: boolean
	created: string
}

export interface IRealmDTO {
	id?: string
	name: string
	realm: string
	isActive: boolean
	notificationChannel: string
	expirationNotice: boolean
	returnNotice: boolean
	locationType: Locations
	needConfirmed: boolean
	hasResponsible: boolean
	hasEmployees: boolean
	hasCommissioningCert: boolean
	hasPreservations: boolean
	hasTransfer: boolean
	verificationSubtractDay: boolean
}
