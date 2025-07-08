export interface IPreservation {
	id: string
	dateStart: number
	dateEnd: number
	notesStart: string
	notesEnd: string
	created: string
}

export interface IPreservationDTO {
	id: string
	instrumentId: string
	dateStart: number
	dateEnd: number
	notesStart: string
	notesEnd: string
}
