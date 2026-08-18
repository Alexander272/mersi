import { FC, useEffect } from 'react'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { useAppSelector } from '@/hooks/redux'
import { getRealm } from '@/features/realms/realmSlice'
import { calcNextVerificationDate } from '@/utils/format'

type Props = {
	data: ICreateFormField
}

export const DateField: FC<Props> = ({ data }) => {
	const { control, setValue, watch } = useFormContext()
	const watchField = watch(data.hide)
	const realm = useAppSelector(getRealm)

	let interval = 0
	let date = ''
	if (data.field == 'nextVerificationDate') {
		interval = watch('instrument.interVerificationInterval')
		date = watch('verification.verificationDate')
	}

	useEffect(() => {
		if (data.field == 'nextVerificationDate') {
			if (interval && date) {
				const newDate = calcNextVerificationDate(date, interval, realm?.verificationSubtractDay)
				setValue('verification.nextVerificationDate', newDate)
			}
		}
	}, [data.field, date, interval, setValue, realm?.verificationSubtractDay])

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
