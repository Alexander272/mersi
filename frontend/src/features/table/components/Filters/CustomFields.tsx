import { FC, useEffect } from 'react'
import { Controller, useFormContext } from 'react-hook-form'
import { CircularProgress, FormControl, InputLabel, MenuItem, Select, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import dayjs from 'dayjs'

import type { IFilter } from '../../types/params'
import { useGetUniqueInstrumentDataQuery } from '../../instrumentApiSlice'
import { DateTextField } from '@/components/DatePicker/DatePicker'

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
						value={dayjs(+field.value * 1000)}
						onChange={value => field.onChange(value?.startOf('d').unix())}
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

const AutocompleteFilter: FC<Props> = ({ index }) => {
	const { control, setValue, watch } = useFormContext<{ filters: IFilter[] }>()
	const field = watch(`filters.${index}.field`)
	const value = watch(`filters.${index}.value`)

	const { data: options, isFetching } = useGetUniqueInstrumentDataQuery(field, { skip: !field })

	useEffect(() => {
		if (options?.data.length && value == '') setValue(`filters.${index}.value`, options?.data[0])
	}, [options, value, setValue, index])

	if (isFetching) return <CircularProgress size={20} />
	return (
		<FormControl fullWidth>
			<InputLabel id={`filters.${index}.value`}>Значение</InputLabel>

			<Controller
				control={control}
				name={`filters.${index}.value`}
				rules={{ required: true }}
				render={({ field, fieldState: { error } }) => (
					<Select
						multiple
						labelId={`filters.${index}.value`}
						value={field.value.split('|')}
						label='Значение'
						error={Boolean(error)}
						onChange={({ target: { value } }) =>
							field.onChange(typeof value === 'string' ? value : value.join('|'))
						}
					>
						{options?.data.map(r => (
							<MenuItem key={r} value={r}>
								{r}
							</MenuItem>
						))}
					</Select>
				)}
			/>
		</FormControl>
	)
}

const ListFilter: FC<Props> = () => {
	return <></>
}
