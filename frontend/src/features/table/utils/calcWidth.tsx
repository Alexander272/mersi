import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import { ColWidth } from '../constants/defaultValues'

export const useCalcWidth = (data: IColumn[]) => {
	let hasFewRows = false
	const width = data.reduce((ac, cur) => {
		if (cur.children && !cur?.hidden) {
			hasFewRows = true
			return ac + cur.children.reduce((ac, cur) => ac + (cur.hidden ? 0 : cur.width || ColWidth), 0)
		}
		return ac + (cur.hidden ? 0 : cur.width || ColWidth)
	}, 12)

	return { width, hasFewRows }
}
