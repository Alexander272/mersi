import { FC, useCallback, useEffect, useState } from 'react'
import { Controller, useFormContext } from 'react-hook-form'
import { CircularProgress, FormControl, InputLabel, MenuItem, Select, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import dayjs from 'dayjs'

import type { IFilter } from '../../types/params'
import { PermRules } from '@/constants/permissions'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useLazyGetDepartmentsQuery } from '@/features/departments/departmentApiSlice'
import { useLazyGetUniqueEmployeeQuery } from '@/features/employees/employeesApiSlice'
import { useGetUniqueInstrumentDataQuery } from '../../instrumentApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getRealm } from '@/features/realms/realmSlice'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { Option, SelectWithFilter } from '@/components/SelectWithFilter/SelectWithFilter'
import { Fallback } from '@/components/Fallback/Fallback'

type Props = {
	index: number
}

export const CustomFields: FC<Props> = ({ index }) => {
	const methods = useFormContext<{ filters: IFilter[] }>()
	const type = methods.watch(`filters.${index}.fieldType`)

	return (
		<>
			{type == 'text' && <StringFilter index={index} />}
			{type == 'number' && <NumberFilter index={index} />}
			{type == 'date' && <DateFilter index={index} />}
			{type == 'short_date' && <ShortDateFilter index={index} />}
			{type == 'list' && <ListFilter index={index} />}
			{type == 'autocomplete' && <AutocompleteFilter index={index} />}
		</>
	)
}

const StringFilter: FC<Props> = ({ index }) => {
	const methods = useFormContext<{ filters: IFilter[] }>()

	return (
		<>
			<FormControl fullWidth sx={{ maxWidth: 170 }}>
				<InputLabel id={`filters.${index}.compareType`}>Условие</InputLabel>
				<Controller
					name={`filters.${index}.compareType`}
					control={methods.control}
					rules={{ required: true }}
					render={({ field, fieldState: { error } }) => (
						<Select
							{...field}
							error={Boolean(error)}
							labelId={`filters.${index}.compareType`}
							label='Условие'
						>
							<MenuItem key='con' value='con'>
								Содержит
							</MenuItem>
							<MenuItem key='like' value='like'>
								Равен
							</MenuItem>
							<MenuItem key='start' value='start'>
								Начинается с
							</MenuItem>
							<MenuItem key='end' value='end'>
								Заканчивается на
							</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

			<TextField
				label='Значение'
				{...methods.register(`filters.${index}.value`, { required: true })}
				error={Boolean(methods.formState.errors?.filters && methods.formState.errors?.filters[index]?.value)}
				fullWidth
			/>
		</>
	)
}
const NumberFilter: FC<Props> = ({ index }) => {
	const methods = useFormContext<{ filters: IFilter[] }>()

	return (
		<>
			<FormControl fullWidth sx={{ maxWidth: 170 }}>
				<InputLabel id={`filters.${index}.compareType`}>Условие</InputLabel>
				<Controller
					name={`filters.${index}.compareType`}
					control={methods.control}
					rules={{ required: true }}
					render={({ field, fieldState: { error } }) => (
						<Select
							{...field}
							error={Boolean(error)}
							labelId={`filters.${index}.compareType`}
							label='Условие'
						>
							<MenuItem key='n_eq' value='eq'>
								Равно
							</MenuItem>
							<MenuItem key='n_gte' value='gte'>
								Больше или равно
							</MenuItem>
							<MenuItem key='n_lte' value='lte'>
								Меньше или равно
							</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

			<TextField
				label='Значение'
				type='number'
				{...methods.register(`filters.${index}.value`, { required: true })}
				error={Boolean(methods.formState.errors?.filters && methods.formState.errors?.filters[index]?.value)}
				fullWidth
			/>
		</>
	)
}
const DateFilter: FC<Props> = ({ index }) => {
	const methods = useFormContext<{ filters: IFilter[] }>()

	return (
		<>
			<FormControl fullWidth sx={{ maxWidth: 170 }}>
				<InputLabel id={`filters.${index}.compareType`}>Условие</InputLabel>
				<Controller
					name={`filters.${index}.compareType`}
					control={methods.control}
					rules={{ required: true }}
					render={({ field, fieldState: { error } }) => (
						<Select
							{...field}
							error={Boolean(error)}
							labelId={`filters.${index}.compareType`}
							label='Условие'
						>
							<MenuItem key='d_eq' value='eq'>
								Равна
							</MenuItem>
							<MenuItem key='d_gte' value='gte'>
								Больше или равна
							</MenuItem>
							<MenuItem key='d_lte' value='lte'>
								Меньше или равна
							</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

			<Controller
				control={methods.control}
				name={`filters.${index}.value`}
				rules={{ required: true }}
				render={({ field, fieldState: { error } }) => (
					<DatePicker
						{...field}
						value={field.value ? dayjs(field.value) : null}
						onChange={value => field.onChange(value?.startOf('d').toISOString())}
						label={'Значение'}
						showDaysOutsideCurrentMonth
						fixedWeekNumber={6}
						slots={{
							textField: DateTextField,
						}}
						slotProps={{
							textField: {
								error: Boolean(error),
							},
						}}
						sx={{ width: '100%' }}
					/>
				)}
			/>
		</>
	)
}

const ShortDateFilter: FC<Props> = ({ index }) => {
	const methods = useFormContext<{ filters: IFilter[] }>()

	return (
		<>
			<FormControl fullWidth sx={{ maxWidth: 170 }}>
				<InputLabel id={`filters.${index}.compareType`}>Условие</InputLabel>
				<Controller
					name={`filters.${index}.compareType`}
					control={methods.control}
					rules={{ required: true }}
					render={({ field, fieldState: { error } }) => (
						<Select
							{...field}
							error={Boolean(error)}
							labelId={`filters.${index}.compareType`}
							label='Условие'
						>
							<MenuItem key='d_eq' value='eq'>
								Равна
							</MenuItem>
							<MenuItem key='d_gte' value='gte'>
								Больше или равна
							</MenuItem>
							<MenuItem key='d_lte' value='lte'>
								Меньше или равна
							</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

			<Controller
				control={methods.control}
				name={`filters.${index}.value`}
				rules={{ required: true }}
				render={({ field, fieldState: { error } }) => (
					<DatePicker
						{...field}
						value={field.value ? dayjs(field.value) : null}
						onChange={value => field.onChange(value?.startOf('m').startOf('d').toISOString())}
						label={'Значение'}
						views={['month', 'year']}
						slots={{
							textField: DateTextField,
						}}
						slotProps={{
							textField: {
								error: Boolean(error),
							},
						}}
						sx={{ width: '100%' }}
					/>
				)}
			/>
		</>
	)
}

const AutocompleteFilter: FC<Props> = ({ index }) => {
	const [options, setOptions] = useState<Option[]>([])

	const section = useAppSelector(getSection)
	const { setValue, watch } = useFormContext<{ filters: IFilter[] }>()
	const field = watch(`filters.${index}.field`)
	const value = watch(`filters.${index}.value`)

	const { data, isFetching } = useGetUniqueInstrumentDataQuery(
		{ field, section: section?.id || '' },
		{ skip: !field || !section?.id }
	)

	useEffect(() => {
		if (!data) return
		setOptions(data.data.map(d => ({ id: d, name: d })))
		if (data.data.length && value == '') setValue(`filters.${index}.value`, data.data[0])
		// if (options?.data.length && value == '') setValue(`filters.${index}.value`, options?.data[0])
	}, [data, value, setValue, index])

	const changeHandler = (values: Option[]) => {
		setValue(`filters.${index}.value`, values.map(v => v.id).join(','))
	}

	if (isFetching) return <CircularProgress size={20} />
	return (
		<SelectWithFilter
			values={options.filter(o => value.includes(o.id))}
			options={options}
			onChange={changeHandler}
		/>
	)
	// return (
	// 	<FormControl fullWidth>
	// 		<InputLabel id={`filters.${index}.value`}>Значение</InputLabel>

	// 		<Controller
	// 			control={control}
	// 			name={`filters.${index}.value`}
	// 			rules={{ required: true }}
	// 			render={({ field, fieldState: { error } }) => (
	// 				<Select
	// 					multiple
	// 					labelId={`filters.${index}.value`}
	// 					value={field.value.split('|')}
	// 					label='Значение'
	// 					error={Boolean(error)}
	// 					onChange={({ target: { value } }) =>
	// 						field.onChange(typeof value === 'string' ? value : value.join('|'))
	// 					}
	// 				>
	// 					{options?.data.map(r => (
	// 						<MenuItem key={r} value={r}>
	// 							{r}
	// 						</MenuItem>
	// 					))}
	// 				</Select>
	// 			)}
	// 		/>
	// 	</FormControl>
	// )
}

export const ListFilter: FC<Props & { label?: string }> = ({ index, label }) => {
	const [options, setOptions] = useState<Option[]>([])

	const hasReserve = useCheckPermission(PermRules.Location.Write)

	const realm = useAppSelector(getRealm)
	const { setValue, watch } = useFormContext<{ filters: IFilter[] }>()
	const field = watch(`filters.${index}.field`)
	const value = watch(`filters.${index}.value`)

	const [getDepartments, { isFetching: isFetchDeps }] = useLazyGetDepartmentsQuery()
	const [getEmployees, { isFetching: isFetchEmp }] = useLazyGetUniqueEmployeeQuery()

	const fetchData = useCallback(async () => {
		if (!realm) return
		if (field == 'place') {
			const { data } = await getDepartments(realm.id)
			const newOptions = [{ id: '_moved', name: 'Перемещение' }]
			if (hasReserve) newOptions.unshift({ id: '_reserve', name: 'Резерв' })
			newOptions.push(...(data?.data.map(d => ({ id: d.id, name: d.name })) || []))
			setOptions(newOptions)
		}
		if (field == 'person') {
			const { data } = await getEmployees(realm.id)
			setOptions(data?.data.map(d => ({
				id: d.id,
				name: d.department ? `${d.name} (${d.department})` : d.name,
			})) || [])
		}
	}, [field, getDepartments, getEmployees, hasReserve, realm])

	useEffect(() => {
		fetchData()
	}, [fetchData])

	const changeHandler = (values: Option[]) => {
		setValue(`filters.${index}.value`, values.map(v => v.id).join(','))
	}

	if (isFetchDeps || isFetchEmp) return <Fallback />
	return (
		<SelectWithFilter
			values={options.filter(o => value.includes(o.id))}
			options={options}
			onChange={changeHandler}
			label={label}
		/>
	)
}
