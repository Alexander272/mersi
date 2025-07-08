import type { IGetSiDTO } from '../types/si'
import { Size } from '@/features/table/constants/defaultValues'

export const buildSiUrlParams = (req: IGetSiDTO): URLSearchParams => {
	const params: string[][] = []

	params.push(['section', req.section])
	if (req.status) params.push(['status', req.status])
	if (req.all) params.push(['all', 'true'])

	if (req.page && req.page != 1) params.push(['page', req.page.toString()])
	if (req.size && req.size != Size) params.push(['size', req.size.toString()])

	if (req.search?.value) {
		params.push([`search[${req.search.fields.join(',')}]`, req.search.value])
	}

	if (req.sort && Object.keys(req.sort).length) {
		console.log(req.sort)
		const sort: string[] = []
		Object.keys(req.sort).forEach(k => {
			sort.push(`${req.sort![k] == 'DESC' ? '-' : ''}${k}`)
		})
		params.push(['sort_by', sort.join(',')])
	}

	if (req.filters) {
		for (let i = 0; i < req.filters.length; i++) {
			const f = req.filters[i]

			if (i == 0 || f.field != req.filters[i - 1].field) {
				params.push([`filters[${f.field}]`, f.fieldType])
			}
			params.push([`${f.field}[${f.compareType}]`, f.value])
		}
	}

	return new URLSearchParams(params)
}
