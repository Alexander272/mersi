import dayjs from 'dayjs'

import type { IFilter } from '../types/params'

const dateFilterFields = new Set(['nextVerificationDate', 'verificationDate', 'dateOfReceipt', 'repairInfo'])

export const normalizeFilters = (filters: IFilter[]): IFilter[] =>
	filters.map(f => {
		if (!dateFilterFields.has(f.field)) return f
		// Старые сохранённые фильтры хранили дату как unix-время в секундах (например 1774983600).
		if (!/^\d{9,11}$/.test(f.value)) return f
		return { ...f, value: dayjs.unix(Number(f.value)).toISOString() }
	})