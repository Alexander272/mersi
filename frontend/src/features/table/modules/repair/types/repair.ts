export interface IRepair {
	id: string
	defect: string
	work: string
	periodStart: number
	periodEnd: number
	description: string
	created: string
}

export interface IRepairDTO {
	id: string
	instrumentId: string
	defect: string
	work: string
	periodStart: number
	periodEnd: number
	description: string
}
