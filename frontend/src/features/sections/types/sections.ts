export interface IGroupedSections {
	id: string
	title: string
	realm: string
	sections: ISection[]
}

export interface ISection {
	id: string
	realmId: string
	name: string
	position: number
	bidType: string
	created: Date
}

export interface ISectionDTO {
	id: string
	name: string
	position: number
	maxPosition: number
	realmId: string
}
