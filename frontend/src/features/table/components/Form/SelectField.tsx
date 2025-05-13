import { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'

type Props = {
	data: ICreateFormField
}

export const SelectField: FC<Props> = ({ data }) => {
	const { control } = useFormContext()

	return (
		<FormControl>
			<InputLabel id={data.field}>{data.fieldName}</InputLabel>
			<Controller
				control={control}
				name={data.path + '.' + data.field}
				render={({ field, fieldState: { error } }) => (
					<Select
						{...field}
						value={field.value || ''}
						labelId={data.field}
						label={data.fieldName}
						error={Boolean(error)}
					>
						{/* //TODO add options */}
						<MenuItem value={''}></MenuItem>
						{/* <MenuItem value={VerificationStatuses.Repair}>Нужен ремонт</MenuItem>
						<MenuItem value={VerificationStatuses.Decommissioning}>Не пригоден</MenuItem> */}
					</Select>
				)}
			/>
		</FormControl>
	)
}
