import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import { useAppSelector } from '@/hooks/redux'
import { ColWidth } from '../constants/defaultValues'
import { getHidden } from '../tableSlice'

export const useCalcWidth = (data: IColumn[]) => {
	const hidden = useAppSelector(getHidden)

	let hasFewRows = false
	const width = data.reduce((ac, cur) => {
		if (cur.children) {
			hasFewRows = true
			return ac + cur.children.reduce((ac, cur) => ac + (hidden[cur.field] ? 0 : cur.width || ColWidth), 0)
		}
		return ac + (hidden[cur.field] ? 0 : cur.width || ColWidth)
	}, 12)

	return { width, hasFewRows }
}
