import dayjs from 'dayjs'

import type { ColumnTypes } from '@/features/sections/modules/columns/types/columns'
import { NullDate } from '@/constants/defaultValues'

const DateFormat = 'DD.MM.YYYY'

type Formatter = (type: ColumnTypes, value: unknown) => string

export const Formatter: Formatter = (type, value) => {
	if (!value || value == NullDate) return '-'

	switch (type) {
		case 'date':
			return dayjs(value as string).format(DateFormat)
		case 'short_date':
			return dayjs(value as string).format('MMMM YYYY')
		// case 'number':
		// 	return new Intl.NumberFormat('ru').format((value as number) || 0)

		default:
			return value as string
	}
}
