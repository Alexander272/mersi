import { useState } from 'react'
import { MenuItem, Select, SelectChangeEvent, Stack, Typography } from '@mui/material'
import { useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IFilter } from '../../types/params'
import { Checkbox } from '@/components/Checkbox/Checkbox'

const months = [
	'Январе',
	'Феврале',
	'Марте',
	'Апреле',
	'Мае',
	'Июне',
	'Июле',
	'Августе',
	'Сентябре',
	'Октябре',
	'Ноябре',
	'Декабре',
]

type ActiveType = 'overdue' | 'month' | 'empty'

export const Default = () => {
	const [active, setActive] = useState<ActiveType | undefined>(
		(localStorage.getItem('active') as ActiveType) || undefined
	)
	const [month, setMonth] = useState(dayjs().get('month'))

	const { setValue, reset } = useFormContext<{ filters: IFilter[] }>()

	const activeHandler = (value: ActiveType) => {
		setActive(prev => (prev != value ? value : undefined))
		if (active != value) localStorage.setItem('active', value)
		else localStorage.removeItem('active')
		reset()
		return active != value
	}

	const emptyHandler = () => {
		const isActive = activeHandler('empty')
		if (!isActive) return

		setValue(`filters.${0}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'eq',
			value: '0',
		})
	}

	const overdueHandler = () => {
		const isActive = activeHandler('overdue')
		if (!isActive) return

		setValue(`filters.${0}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'lte',
			value: dayjs().startOf('d').unix().toString(),
		})
	}

	const monthHandler = () => {
		const isActive = activeHandler('month')
		if (!isActive) return

		const date = dayjs().set('month', month)
		setValue(`filters.${0}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'gte',
			value: date.startOf('month').unix().toString(),
		})
		setValue(`filters.${1}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'lte',
			value: date.endOf('month').unix().toString(),
		})
	}
	const curMonthHandler = (event: SelectChangeEvent<number>) => {
		setMonth(+event.target.value)
		reset()
		const date = dayjs().set('month', +event.target.value)
		setValue(`filters.${0}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'gte',
			value: date.startOf('month').unix().toString(),
		})
		setValue(`filters.${1}`, {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'lte',
			value: date.endOf('month').unix().toString(),
		})
	}

	return (
		<Stack spacing={1}>
			<Typography>Показать все инструменты у которых:</Typography>

			<Checkbox
				id='empty'
				name='empty'
				checked={active == 'empty'}
				onChange={emptyHandler}
				label='Срок следующей поверки не задан'
			/>
			<Checkbox
				id='overdue'
				name='overdue'
				checked={active == 'overdue'}
				onChange={overdueHandler}
				label='Срок следующей поверки прошел'
			/>

			<Stack direction={'row'} justifyContent={'space-between'}>
				<Checkbox
					id='month'
					name='month'
					checked={active == 'month'}
					onChange={monthHandler}
					label='Следующая поверка в'
				/>
				<Select value={month} onChange={curMonthHandler} disabled={active != 'month'} sx={{ width: 250 }}>
					{months.map((m, i) => (
						<MenuItem key={m} value={i}>
							{m}
						</MenuItem>
					))}
				</Select>
			</Stack>
		</Stack>
	)
}
