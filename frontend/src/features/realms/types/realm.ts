type Locations = 'department' | 'place'

export interface IRealm {
	id: string
	name: string
	realm: string
	isActive: boolean
	reserveChannel: string
	expirationNotice: boolean
	locationType: Locations
	hasResponsible: boolean
	needResponsible: boolean
	needConfirmed: boolean
	created: string
}

export interface IRealmDTO {
	id?: string
	name: string
	realm: string
	isActive: boolean
	reserveChannel: string
	expirationNotice: boolean
	locationType: Locations
	hasResponsible: boolean
	needResponsible: boolean
	needConfirmed: boolean
}
