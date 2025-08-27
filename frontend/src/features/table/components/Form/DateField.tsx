import { FC, useEffect } from 'react'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { DateTextField } from '@/components/DatePicker/DatePicker'

type Props = {
	data: ICreateFormField
}

export const DateField: FC<Props> = ({ data }) => {
	const { control, setValue, watch } = useFormContext()
	const watchField = watch(data.hide)

	let interval = 0
	let date = 0
	if (data.field == 'nextVerificationDate') {
		interval = watch('instrument.interVerificationInterval')
		date = watch('verification.verificationDate')
	}

	useEffect(() => {
		if (data.field == 'nextVerificationDate') {
			if (interval && date) {
				const newDate = dayjs(date * 1000)
					.add(interval, 'M')
					.subtract(1, 'd')
					.unix()
				setValue('verification.nextVerificationDate', newDate)
			}
		}
	}, [data.field, date, interval, setValue])

	if (data.hide && watchField) return null
	return (
		<Controller
			control={control}
			name={data.path + '.' + data.field}
			rules={{ required: true, min: 1000000000 }}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={dayjs(field.value * 1000)}
					onChange={value => {
						field.onChange(value?.startOf('d').unix())
					}}
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
