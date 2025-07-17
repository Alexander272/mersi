export interface IWriteOff {
	id: string
	instrumentId: string
	date: number
	notes: string
	docId: string
	docName: string
	created: string
}

export interface IWriteOffDTO {
	id: string
	instrumentId: string
	date: number
	notes: string
	docId: string
	docName: string
}
