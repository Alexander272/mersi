import dayjs from 'dayjs'

import { NullDate } from '@/constants/defaultValues'

export const numberFormat = new Intl.NumberFormat('ru-Ru').format

export const isEmptyDate = (v?: string) => !v || v == NullDate

export const calcNextVerificationDate = (date: string, interval: number, subtractDay?: boolean) => {
	const next = dayjs(date).add(interval, 'month')
	return subtractDay ? next.subtract(1, 'd').toISOString() : next.toISOString()
}

export const removeSpace = (values: string[]) => {
	return values.map(v => v.replace(/\s+/g, ''))
}
