import { FC, useEffect } from 'react'
import { Autocomplete, TextField } from '@mui/material'
import { Controller, RegisterOptions, useFormContext } from 'react-hook-form'

import type { ILocationForm } from '../../../types/location'
import { useGetEmployeesQuery } from '@/features/employees/employeesApiSlice'

type Props = {
	label: string
	name: 'person'
	rules?: RegisterOptions<ILocationForm, 'person'>
	disabled?: boolean
}

export const EmployeeList: FC<Props> = ({ label, name, rules, disabled }) => {
	const { control, watch, setValue } = useFormContext<ILocationForm>()

	const departmentId = watch('department')
	const { data, isFetching } = useGetEmployeesQuery(departmentId || null, { skip: !departmentId })

	useEffect(() => {
		if (data?.data.length) setValue(name, data.data[0].id)
		else setValue(name, '')
	}, [name, data, setValue])

	return (
		<Controller
			control={control}
			name={name}
			rules={rules}
			disabled={disabled}
			render={({ field, fieldState: { error } }) => (
				<Autocomplete
					value={data?.data.find(d => d.id == field.value) || ''}
					onChange={(_event, value) => {
						field.onChange(typeof value == 'string' ? value : value.id)
					}}
					options={data?.data || []}
					loading={isFetching}
					getOptionLabel={option => (typeof option === 'string' ? option : option.name)}
					freeSolo
					disableClearable
					noOptionsText='Ничего не найдено'
					loadingText='Загрузка...'
					disabled={disabled}
					fullWidth
					renderInput={params => (
						<TextField {...params} label={label} autoComplete='off' error={Boolean(error)} />
					)}
				/>
			)}
		/>
	)
}
