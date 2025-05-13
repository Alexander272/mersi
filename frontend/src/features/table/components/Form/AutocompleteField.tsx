import { FC } from 'react'
import { Autocomplete, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { useGetUniqueInstrumentDataQuery } from '../../instrumentApiSlice'

type Props = {
	data: ICreateFormField
}

export const AutocompleteField: FC<Props> = ({ data }) => {
	const { control } = useFormContext()

	const { data: options, isFetching } = useGetUniqueInstrumentDataQuery(data.field, { skip: !data.field })

	return (
		<Controller
			name={data.path + '.' + data.field}
			control={control}
			rules={{ required: data.isRequired }}
			render={({ field: { onChange, value, ref }, fieldState: { error } }) => (
				<Autocomplete
					value={value || ''}
					freeSolo
					disableClearable
					autoComplete
					options={options?.data || []}
					loading={isFetching}
					// TODO maybe I have to add icon
					onChange={(_event, value) => {
						onChange(value)
					}}
					renderInput={params => (
						<TextField
							{...params}
							label={data.fieldName}
							onChange={onChange}
							error={Boolean(error)}
							helperText={error?.message}
							inputRef={ref}
						/>
					)}
				/>
			)}
		/>
	)
}
