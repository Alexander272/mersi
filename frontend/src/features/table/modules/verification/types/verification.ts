import { IVerDocs } from './verificationDocs'

export interface IVerification {
	id: string
	verificationDate: string
	nextVerificationDate: string
	registerLink: string
	notVerified: boolean
	status: string
	notes: string
	docs?: IVerDocs[]
}

export interface IVerificationDTO {
	id: string
	instrumentId: string
	verificationDate: string
	nextVerificationDate: string
	registerLink: string
	notVerified: boolean
	status: string
	notes: string
	docs?: IVerDocs[]
}
