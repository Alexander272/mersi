import { FC } from 'react'
import { FormControl, TextField as MTextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'

type Props = {
	data: ICreateFormField
}

export const TextField: FC<Props> = ({ data }) => {
	const { control } = useFormContext()

	return (
		<FormControl>
			<Controller
				control={control}
				name={data.path + '.' + data.field}
				rules={{ required: data.isRequired }}
				render={({ field, fieldState: { error } }) => (
					<MTextField
						{...field}
						value={field.value || ''}
						label={data.fieldName}
						fullWidth
						error={Boolean(error)}
						// multiline={NotesField.multiline}
						// minRows={NotesField.minRows}
					/>
				)}
			/>
		</FormControl>
	)
}
