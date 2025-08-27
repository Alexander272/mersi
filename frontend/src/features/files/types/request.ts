export type IGetDocumentsDTO = {
	group: string
	instrument: string
}

export interface IUploadFiles {
	data: FormData
}

export type DeleteDocuments = {
	id: string
	filename: string
	group: string
	instrumentId: string
	isTemp: boolean
}
