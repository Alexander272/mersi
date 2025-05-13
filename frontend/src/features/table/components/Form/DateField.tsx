import { FC } from 'react'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { DateTextField } from '@/components/DatePicker/DatePicker'

type Props = {
	data: ICreateFormField
}

export const DateField: FC<Props> = ({ data }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={data.path + '.' +data.field}
			rules={{ required: true, min: 1000000000 }}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={dayjs(field.value * 1000)}
					onChange={value => field.onChange(value?.startOf('d').unix())}
					label={data.fieldName}
					showDaysOutsideCurrentMonth
					fixedWeekNumber={6}
					// disableFuture
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
