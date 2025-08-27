import { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { VerificationStatuses } from '../../modules/verification/constants/status'
import { DepartmentList } from '../../modules/locations/components/Forms/SendToPerson/DepartmentList'
import { EmployeeList } from '../../modules/locations/components/Forms/SendToPerson/EmployeeList'

type Props = {
	data: ICreateFormField
}

export const SelectField: FC<Props> = ({ data }) => {
	const { control, watch } = useFormContext()
	const watchField = watch(data.hide)

	if (data.hide && watchField) return null
	if (data.field == 'department') return <DepartmentList label={data.fieldName} name={data.field} />
	if (data.field == 'person') return <EmployeeList label={data.fieldName} name={data.field} />
	return (
		<FormControl>
			<InputLabel id={data.field}>{data.fieldName}</InputLabel>
			<Controller
				control={control}
				name={`${data.path}.${data.field}`}
				render={({ field, fieldState: { error } }) => (
					<Select
						{...field}
						value={field.value || ''}
						labelId={data.field}
						label={data.fieldName}
						error={Boolean(error)}
					>
						<MenuItem value={VerificationStatuses.Work}>Пригоден</MenuItem>
						<MenuItem value={VerificationStatuses.Repair}>Нужен ремонт</MenuItem>
						<MenuItem value={VerificationStatuses.Decommissioning}>Не пригоден</MenuItem>
					</Select>
				)}
			/>
		</FormControl>
	)
}
