import dayjs from 'dayjs'

import type { ColumnTypes } from '@/features/sections/modules/columns/types/columns'
import { NullDate } from '@/constants/defaultValues'

const DateFormat = 'DD.MM.YYYY'
const ShortDateFormat = 'MMMM YYYY'

type Formatter = (type: ColumnTypes, value: unknown) => string

export const Formatter: Formatter = (type, value) => {
	if (value == null || value === '' || value === NullDate) return '-'

	if (type === 'date' || type === 'short_date') {
		const strValue = String(value)
		if (strValue.startsWith('0001-01-01')) return '-'

		const date = dayjs(value as string)
		if (!date.isValid()) return '-'

		return date.format(type === 'date' ? DateFormat : ShortDateFormat)
	}

	if (type == 'number' && value == 0) return '-'

	// if (type === 'number') return new Intl.NumberFormat('ru').format(Number(value) || 0)

	return String(value)
}
