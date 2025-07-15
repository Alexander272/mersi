import { FC } from 'react'
import { Stack, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import { DateTextField } from '@/components/DatePicker/DatePicker'

type Props = {
	isThisSave: boolean
	min?: number
}

export const Inputs: FC<Props> = ({ isThisSave, min = 1262286000 }) => {
	const { control } = useFormContext()

	return (
		<Stack spacing={2} mb={2}>
			<Controller
				control={control}
				name={isThisSave ? 'dateStart' : 'dateEnd'}
				render={({ field, fieldState: { error } }) => (
					<DatePicker
						{...field}
						value={dayjs(field.value * 1000)}
						onChange={value => field.onChange(value?.startOf('d').unix())}
						label={`Дата ${isThisSave ? 'передачи' : 'возврата'}`}
						showDaysOutsideCurrentMonth
						fixedWeekNumber={6}
						minDate={dayjs(min * 1000)}
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

			<Controller
				control={control}
				name={isThisSave ? 'notesStart' : 'notesEnd'}
				render={({ field }) => <TextField {...field} label={'Примечание'} multiline minRows={4} />}
			/>
		</Stack>
	)
}
