import { IInstrumentDTO } from './instrument'
import { IVerificationDTO } from './verification'

export interface ISiDTO {
	instrument: IInstrumentDTO
	verification?: IVerificationDTO
}
