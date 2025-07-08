import { FC, useEffect, useRef } from 'react'
import { FormControl, IconButton, InputLabel, MenuItem, Select, SelectChangeEvent, Stack } from '@mui/material'
import { Controller, UseFieldArrayReturn, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { CompareTypes, IFilter } from '../../types/params'
import { useAppSelector } from '@/hooks/redux'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { CustomFields } from './CustomFields'

type Props = {
	methods: UseFieldArrayReturn<{ filters: IFilter[] }, 'filters', 'id'>
}

export const Custom: FC<Props> = ({ methods }) => {
	const { fields, append, remove } = methods
	const isEmpty = useRef(fields.length == 0)

	const removeHandler = (index: number) => {
		remove(index)
	}

	useEffect(() => {
		if (isEmpty.current) {
			isEmpty.current = false
			append({
				field: 'name',
				fieldType: 'text',
				compareType: 'con' as const,
				value: '',
			})
		}
	}, [append, fields])

	return (
		<Stack spacing={2}>
			{fields.map((f, i) => (
				<FilterItem key={f.id} index={i} remove={removeHandler} />
			))}
		</Stack>
	)
}

const compareTypes = new Map([
	['text', 'con'],
	['date', 'eq'],
	['number', 'eq'],
	// ['switch', 'eq'],
	['list', 'in'],
	['autocomplete', 'in'],
])
const defaultValues = new Map([
	['text', ''],
	['date', dayjs().unix().toString()],
	['number', ''],
	// ['switch', 'false'],
])

type FilterItemProps = {
	index: number
	remove: (index: number) => void
}
const FilterItem: FC<FilterItemProps> = ({ index, remove }) => {
	const section = useAppSelector(getSection)

	const methods = useFormContext<{ filters: IFilter[] }>()
	const type = methods.watch(`filters.${index}.fieldType`)

	const { data } = useGetColumnsQuery(section?.id || '', { skip: !section?.id })

	const removeHandler = () => remove(index)

	return (
		<Stack direction={'row'} spacing={1} alignItems={'center'}>
			<Controller
				control={methods.control}
				name={`filters.${index}.field`}
				render={({ field, fieldState: { error } }) => (
					<FormControl fullWidth sx={{ maxWidth: 170 }}>
						<InputLabel id={`filters.${index}.field`}>Колонка</InputLabel>

						<Select
							value={`${field.value}@${type}`}
							onChange={(event: SelectChangeEvent) => {
								const newType = event.target.value.split('@')[1]
								if (newType != type) {
									methods.setValue(`filters.${index}.fieldType`, newType as 'text')
									methods.setValue(
										`filters.${index}.compareType`,
										(compareTypes.get(newType) || 'con') as CompareTypes
									)
								}
								if (newType != type || newType == 'autocomplete' || newType == 'list') {
									methods.setValue(`filters.${index}.value`, defaultValues.get(newType) || '')
								}
								field.onChange(event.target.value.split('@')[0])
							}}
							labelId={`filters.${index}.field`}
							label={'Колонка'}
							error={Boolean(error)}
						>
							{data?.data.map(c => (
								<MenuItem key={c.field} value={`${c.field}@${c.type}`}>
									{c.name}
								</MenuItem>
							))}
						</Select>
					</FormControl>
				)}
			/>

			<CustomFields index={index} />

			{index != 0 && (
				<IconButton onClick={removeHandler}>
					<TimesIcon fontSize={18} padding={0.4} />
				</IconButton>
			)}
		</Stack>
	)
}
