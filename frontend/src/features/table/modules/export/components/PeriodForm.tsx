import { FC } from 'react'
import { Button, Stack } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, FormProvider, useForm } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IPeriod } from '../types/period'
import { getSection } from '@/features/sections/sectionSlice'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useLazyMakeSchedulerQuery } from '../exportApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { DateTextField } from '@/components/DatePicker/DatePicker'

const defaultValues: IPeriod = {
	gte: dayjs().startOf('D').startOf('M').toISOString(),
	lte: dayjs().add(1, 'M').startOf('D').startOf('M').toISOString(),
	section: '',
}
const min = '2000-01-01'

export const PeriodForm: FC = () => {
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const methods = useForm<IPeriod>({ values: defaultValues })

	const [makeSchedule] = useLazyMakeSchedulerQuery()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Schedule', isOpen: false }))
	}

	const submitHandler = methods.handleSubmit(async data => {
		if (!section) return
		await makeSchedule({ ...data, section: section?.id || '' }).unwrap()
		closeHandler()
	})

	return (
		<Stack component={'form'} onSubmit={submitHandler} mt={2}>
			<FormProvider {...methods}>
				<Stack spacing={2}>
					<Controller
						control={methods.control}
						name={'gte'}
						rules={{ required: true }}
						render={({ field, fieldState: { error } }) => (
							<DatePicker
								{...field}
								value={dayjs(field.value)}
								onChange={value => field.onChange(value?.startOf('d').toISOString())}
								label={'Начало периода'}
								showDaysOutsideCurrentMonth
								fixedWeekNumber={6}
								minDate={dayjs(min)}
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
						control={methods.control}
						name={'lte'}
						rules={{ required: true }}
						render={({ field, fieldState: { error } }) => (
							<DatePicker
								{...field}
								value={dayjs(field.value)}
								onChange={value => field.onChange(value?.startOf('d').toISOString())}
								label={'Конец периода'}
								showDaysOutsideCurrentMonth
								fixedWeekNumber={6}
								minDate={dayjs(min)}
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
				</Stack>
			</FormProvider>

			<Stack direction={'row'} spacing={3} mt={4}>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отменить
				</Button>
				<Button variant='contained' type='submit' fullWidth>
					Применить
				</Button>
			</Stack>
		</Stack>
	)
}
