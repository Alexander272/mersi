import { FC } from 'react'
import { FormControl, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'

type Props = {
	data: ICreateFormField
}

export const NumberField: FC<Props> = ({ data }) => {
	const { control } = useFormContext()

	return (
		<FormControl>
			<Controller
				control={control}
				name={data.path + '.' + data.field}
				rules={{ required: data.isRequired }}
				render={({ field, fieldState: { error } }) => (
					<TextField
						{...field}
						value={field.value || ''}
						onChange={e => field.onChange(+(e.target.value || 0))}
						label={data.fieldName}
						fullWidth
						error={Boolean(error)}
						slotProps={{
							htmlInput: {
								type: 'number',
								step: 1,
								// min: 1,
								// max: 100
							},
						}}
					/>
				)}
			/>
		</FormControl>
	)
}
