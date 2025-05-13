export type GetDocuments = {
	verificationId: string
	instrumentId: string
}

export interface IUploadFiles {
	data: FormData
}

export type DeleteDocuments = {
	id: string
	filename: string
	group: string
	instrumentId: string
}
