import { FC } from 'react'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { DateTextField } from '@/components/DatePicker/DatePicker'

type Props = {
	data: ICreateFormField
}

export const ShortDateField: FC<Props> = ({ data }) => {
	const { control, watch } = useFormContext()
	const watchField = watch(data.hide)

	if (data.hide && watchField) return null
	return (
		<Controller
			control={control}
			name={data.path + '.' + data.field}
			rules={{ required: data.isRequired }}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={field.value ? dayjs(field.value) : null}
					onChange={value => {
						field.onChange(value?.startOf('d').toISOString())
					}}
					label={data.fieldName}
					views={['month', 'year']}
					disableFuture
					slots={{
						textField: DateTextField,
					}}
					slotProps={{
						textField: {
							error: Boolean(error),
						},
					}}
				/>
			)}
		/>
	)
}
