import { FC } from 'react'
import { Checkbox, FormControl, FormControlLabel, useTheme } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'

type Props = {
	data: ICreateFormField
}

export const CheckboxField: FC<Props> = ({ data }) => {
	const { palette } = useTheme()
	const { control, watch } = useFormContext()
	const watchField = watch(data.hide)

	if (data.hide && watchField) return null
	return (
		<FormControl>
			<Controller
				control={control}
				name={data.path + '.' + data.field}
				rules={{ required: data.isRequired }}
				render={({ field }) => (
					<FormControlLabel
						label={data.fieldName}
						control={<Checkbox checked={field.value || false} />}
						onChange={field.onChange}
						sx={{
							transition: 'all 0.3s ease-in-out',
							borderRadius: 3,
							userSelect: 'none',
							margin: 0,
							':hover': { backgroundColor: palette.action.hover },
						}}
					/>
				)}
			/>
		</FormControl>
	)
}
