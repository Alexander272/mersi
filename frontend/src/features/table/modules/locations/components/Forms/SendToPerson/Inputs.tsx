import { Stack } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { ILocationForm } from '../../../types/location'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { Checkbox } from '@/components/Checkbox/Checkbox'
import { DepartmentList } from './DepartmentList'
import { EmployeeList } from './EmployeeList'
import { FC } from 'react'
import { useAppSelector } from '@/hooks/redux'
import { getRealm } from '@/features/realms/realmSlice'

type Props = {
	minDate?: number
}

export const Inputs: FC<Props> = ({ minDate = 1262286000 }) => {
	const { control } = useFormContext<ILocationForm>()
	const realm = useAppSelector(getRealm)

	return (
		<Stack spacing={2} alignItems={'flex-start'}>
			{/* <Controller
				name='isToReserve'
				control={control}
				render={({ field }) => (
					<Checkbox id='isToReserve' label='В резерве' {...field} checked={field.value} value={''} />
				)}
			/> */}
			{realm?.needConfirmed && (
				<Controller
					name='needConfirm'
					control={control}
					render={({ field }) => (
						<Checkbox
							id='needConfirm'
							label='Нужны уведомления'
							{...field}
							checked={field.value}
							value={''}
						/>
					)}
				/>
			)}

			<DepartmentList label={'Подразделение'} name='department' rules={{ required: true }} />
			{realm?.hasResponsible && (
				<EmployeeList label={'Лицо держащее СИ'} name='person' rules={{ required: realm.needResponsible }} />
			)}

			<Controller
				control={control}
				name={'dateOfIssue'}
				render={({ field, fieldState: { error } }) => (
					<DatePicker
						{...field}
						value={dayjs(field.value * 1000)}
						onChange={value => field.onChange(value?.startOf('d').unix())}
						label={'Дата выдачи или поступления'}
						showDaysOutsideCurrentMonth
						fixedWeekNumber={6}
						minDate={dayjs(minDate * 1000)}
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
	)
}
