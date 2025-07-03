import { IVerDocs } from './verificationDocs'

export interface IVerificationDTO {
	id: string
	instrumentId: string
	verificationDate: number
	nextVerificationDate: number
	registerLink: string
	notVerified: boolean
	status: string
	notes: string
	docs?: IVerDocs[]
}
