export interface IRepair {
	id: string
	instrumentId: string
	defect: string
	work: string
	periodStart: string
	periodEnd: string
	description: string
	created: string
}

export interface IRepairDTO {
	id: string
	instrumentId: string
	defect: string
	work: string
	periodStart: string
	periodEnd: string
	description: string
}
