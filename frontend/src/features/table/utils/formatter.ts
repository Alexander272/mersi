import dayjs from 'dayjs'

import type { ColumnTypes } from '@/features/sections/modules/columns/types/columns'

const DateFormat = 'DD.MM.YYYY'

type Formatter = (type: ColumnTypes, value: unknown) => string

export const Formatter: Formatter = (type, value) => {
	if (!value) return '-'

	switch (type) {
		case 'date':
			return dayjs((value as number) * 1000).format(DateFormat)
		// case 'number':
		// 	return new Intl.NumberFormat('ru').format((value as number) || 0)

		default:
			return value as string
	}
}
