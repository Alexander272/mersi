import { useEffect, useState } from 'react'
import { MenuItem, Select, SelectChangeEvent, Stack, Typography } from '@mui/material'
import { useFieldArray, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IFilter } from '../../types/params'
import { localKeys } from '@/constants/localKeys'
import { PermRules } from '@/constants/permissions'
import { NullDate } from '@/constants/defaultValues'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { useGetDepartmentsQuery } from '@/features/departments/departmentApiSlice'
import { useGetUniqueEmployeeQuery } from '@/features/employees/employeesApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getRealm } from '@/features/realms/realmSlice'
import { getFilters } from '../../tableSlice'
import { SelectWithFilter, type Option } from '@/components/SelectWithFilter/SelectWithFilter'
import { Checkbox } from '@/components/Checkbox/Checkbox'
import { Fallback } from '@/components/Fallback/Fallback'

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

type ActiveType = 'overdue' | 'month' | 'empty' | 'place' | 'person'

export const Default = () => {
	const filters = useAppSelector(getFilters)

	const [active, setActive] = useState<ActiveType | undefined>(
		(localStorage.getItem(localKeys.activeFilters) as ActiveType) || undefined
	)
	const [month, setMonth] = useState(active == 'month' ? dayjs(filters[0]?.value).get('month') : dayjs().get('month'))
	const section = useAppSelector(getSection)

	const { control } = useFormContext<{ filters: IFilter[] }>()
	const { replace } = useFieldArray({ control, name: 'filters' })

	const { data } = useGetColumnsQuery({ section: section?.id || '', original: true }, { skip: !section?.id })

	const activeHandler = (value: ActiveType) => {
		if (active != value) {
			setActive(value)
			localStorage.setItem(localKeys.activeFilters, value)
		} else {
			setActive(undefined)
			localStorage.removeItem(localKeys.activeFilters)
		}
		replace([])
		return active != value
	}

	const emptyHandler = () => {
		const isActive = activeHandler('empty')
		if (!isActive) return

		replace([{ field: 'nextVerificationDate', fieldType: 'date', compareType: 'eq', value: NullDate }])
	}

	const overdueHandler = () => {
		const isActive = activeHandler('overdue')
		if (!isActive) return

		replace([
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'lte',
				value: dayjs().startOf('d').toISOString(),
			},
		])
	}

	const monthHandler = () => {
		const isActive = activeHandler('month')
		if (!isActive) return

		const date = dayjs().set('month', month)
		replace([
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'gte',
				value: date.startOf('month').toISOString(),
			},
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'lte',
				value: date.endOf('month').toISOString(),
			},
		])
	}
	const curMonthHandler = (event: SelectChangeEvent<number>) => {
		setMonth(+event.target.value)
		const date = dayjs().set('month', +event.target.value)
		replace([
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'gte',
				value: date.startOf('month').toISOString(),
			},
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'lte',
				value: date.endOf('month').toISOString(),
			},
		])
	}

	const placeHandler = () => {
		const isActive = activeHandler('place')
		if (!isActive) return
		replace([{ field: 'place', fieldType: 'list', compareType: 'in', value: '' }])
	}

	const personHandler = () => {
		const isActive = activeHandler('person')
		if (!isActive) return
		replace([{ field: 'person', fieldType: 'list', compareType: 'in', value: '' }])
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

			{data?.data.find(col => col.field == 'place') && (
				<Stack direction={'row'} justifyContent={'space-between'}>
					<Checkbox
						id='place'
						name='place'
						checked={active == 'place'}
						onChange={placeHandler}
						label='Место нахождения'
					/>
					<DepartmentFilter disabled={active != 'place'} />
				</Stack>
			)}

			{data?.data.find(col => col.field == 'person') && (
				<Stack direction={'row'} justifyContent={'space-between'}>
					<Checkbox
						id='person'
						name='person'
						checked={active == 'person'}
						onChange={personHandler}
						label='ФИО сотрудника'
					/>
					<PersonFilter disabled={active != 'person'} />
				</Stack>
			)}
		</Stack>
	)
}

type ListProps = { disabled?: boolean }
const DepartmentFilter = ({ disabled }: ListProps) => {
	const [options, setOptions] = useState<Option[]>([])

	const hasReserve = useCheckPermission(PermRules.Location.Write)

	const realm = useAppSelector(getRealm)

	const { setValue, watch } = useFormContext<{ filters: IFilter[] }>()
	const value = watch(`filters.0.value`)

	const { data, isFetching } = useGetDepartmentsQuery(realm?.id || '', { skip: !realm })

	useEffect(() => {
		if (!data) return
		const newOptions = [{ id: '_moved', name: 'Перемещение' }]
		if (hasReserve) newOptions.unshift({ id: '_reserve', name: 'Резерв' })
		newOptions.push(...(data?.data.map(d => ({ id: d.id, name: d.name })) || []))
		setOptions(newOptions)
	}, [data, hasReserve])

	const changeHandler = (values: Option[]) => {
		setValue(`filters.0.value`, values.map(v => v.id).join(','))
	}

	if (isFetching) return <Fallback />
	return (
		<SelectWithFilter
			values={options.filter(o => value?.includes(o.id))}
			options={options}
			onChange={changeHandler}
			disabled={disabled}
			sx={{ width: 270 }}
		/>
	)
}
const PersonFilter = ({ disabled }: ListProps) => {
	const [options, setOptions] = useState<Option[]>([])

	const realm = useAppSelector(getRealm)

	const { data, isFetching } = useGetUniqueEmployeeQuery(realm?.id || '', { skip: !realm })

	const { setValue, watch } = useFormContext()
	const value = watch(`filters.0.value`)

	useEffect(() => {
		if (!data) return
		setOptions(data.data.map(d => ({
			id: d.id,
			name: d.department ? `${d.name} (${d.department})` : d.name,
		})))
	}, [data])

	const changeHandler = (values: Option[]) => {
		setValue(`filters.0.value`, values.map(v => v.id).join(','))
	}

	if (isFetching) return <Fallback />
	return (
		<SelectWithFilter
			values={options.filter(o => value?.includes(o.id))}
			options={options}
			onChange={changeHandler}
			disabled={disabled}
			sx={{ width: 270 }}
		/>
	)
}
