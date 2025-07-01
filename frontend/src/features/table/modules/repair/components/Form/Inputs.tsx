import { Controller, useFormContext } from 'react-hook-form'
import { Stack, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import dayjs from 'dayjs'

import type { IRepairDTO } from '../../types/repair'
import { DateTextField } from '@/components/DatePicker/DatePicker'

export const Inputs = () => {
	const { control } = useFormContext<IRepairDTO>()

	return (
		<Stack spacing={2} mb={2}>
			<Controller
				control={control}
				name='defect'
				render={({ field }) => <TextField {...field} value={field.value || ''} fullWidth label='Дефект' />}
			/>
			<Controller
				control={control}
				name='work'
				render={({ field }) => (
					<TextField {...field} value={field.value || ''} fullWidth label='Проведенные работы' />
				)}
			/>

			<Stack direction={'row'} spacing={2}>
				<Controller
					control={control}
					name='periodStart'
					render={({ field, fieldState: { error } }) => (
						<DatePicker
							{...field}
							value={dayjs(field.value * 1000)}
							onChange={value => field.onChange(value?.startOf('d').unix())}
							label={'Начало ремонта'}
							showDaysOutsideCurrentMonth
							fixedWeekNumber={6}
							slots={{
								textField: DateTextField,
							}}
							slotProps={{
								textField: {
									error: Boolean(error),
								},
							}}
							sx={{ width: '100%' }}
						/>
					)}
				/>

				<Controller
					control={control}
					name='periodEnd'
					render={({ field, fieldState: { error } }) => (
						<DatePicker
							{...field}
							value={dayjs(field.value * 1000)}
							onChange={value => field.onChange(value?.startOf('d').unix())}
							label={'Конец ремонта'}
							showDaysOutsideCurrentMonth
							fixedWeekNumber={6}
							slots={{
								textField: DateTextField,
							}}
							slotProps={{
								textField: {
									error: Boolean(error),
								},
							}}
							sx={{ width: '100%' }}
						/>
					)}
				/>
			</Stack>

			<Controller
				control={control}
				name='description'
				render={({ field }) => (
					<TextField {...field} value={field.value || ''} fullWidth label='Комментарий' multiline rows={3} />
				)}
			/>
		</Stack>
	)
}
