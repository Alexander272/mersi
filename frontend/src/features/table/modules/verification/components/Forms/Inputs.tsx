import { FC, useState } from 'react'
import { FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IDocument } from '@/features/files/types/document'
import { VerificationStatuses } from '../../constants/status'
import { useAppSelector } from '@/hooks/redux'
import { useGetVerificationFieldsQuery } from '../../verificationApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { UploadButton } from '@/features/files/components/UploadButton/UploadButton'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

type Props = {
	instrumentId: string
}

export const Inputs: FC<Props> = ({ instrumentId }) => {
	const section = useAppSelector(getSection)

	const { data, isFetching } = useGetVerificationFieldsQuery(section?.id || '', { skip: !section?.id })

	return (
		<Stack spacing={2} mb={2}>
			{isFetching && <BoxFallback />}

			{data?.data.map(f => {
				if (f.type == 'date') return <DateField key={f.id} field={f.field} label={f.label} />
				if (f.field == 'status') return <StatusField key={f.id} field={f.field} label={f.label} />
				if (f.field == 'registerLink') return <LinkField key={f.id} field={f.field} label={f.label} />
				if (f.field == 'notes') return <NotesField key={f.id} field={f.field} label={f.label} />
				if (f.type == 'file')
					return <FileField key={f.id} field={f.field} label={f.label} instrumentId={instrumentId} />
			})}
		</Stack>
	)
}

type FieldProps = {
	label: string
	field: string
	instrumentId?: string
}

const DateField: FC<FieldProps> = ({ label, field }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={field}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={dayjs(field.value * 1000)}
					onChange={value => field.onChange(value?.startOf('d').unix())}
					label={label}
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
	)
}

const StatusField: FC<FieldProps> = ({ label, field }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={field}
			render={({ field, fieldState: { error } }) => (
				<FormControl>
					<InputLabel id={'status'}>{label}</InputLabel>

					<Select labelId={'status'} label={label} error={Boolean(error)} {...field}>
						<MenuItem value={VerificationStatuses.Work}>Пригоден</MenuItem>
						<MenuItem value={VerificationStatuses.Repair}>Нужен ремонт</MenuItem>
						<MenuItem value={VerificationStatuses.Decommissioning}>Не пригоден</MenuItem>
					</Select>
				</FormControl>
			)}
		/>
	)
}

const LinkField: FC<FieldProps> = ({ label, field }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={field}
			render={({ field }) => <TextField {...field} label={label} multiline />}
		/>
	)
}

const NotesField: FC<FieldProps> = ({ label, field }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={field}
			render={({ field }) => <TextField {...field} label={label} multiline minRows={4} />}
		/>
	)
}

const FileField: FC<FieldProps> = ({ label, field, instrumentId = '' }) => {
	const [doc, setDoc] = useState<IDocument | null>(null)

	const { control, setValue } = useFormContext()

	//TODO надо еще вставлять название файла в поле для ввода и решить что делать если файла не будет, а будет только название
	// еще наверное надо как-то редактировать название файла при изменении его в поле ввода

	const setDocument = (value: IDocument | null) => {
		setDoc(value)
		setValue(field, value?.label || '')
		setValue(`${field}Id`, value?.id || '')
	}

	return (
		<Stack direction={'row'}>
			<Controller
				control={control}
				name={field}
				render={({ field, fieldState: { error } }) => (
					<TextField
						{...field}
						value={field.value || ''}
						label={label}
						error={Boolean(error)}
						sx={{
							flexGrow: 1,
							'.MuiInputBase-root': { borderTopRightRadius: 0, borderBottomRightRadius: 0 },
						}}
					/>
				)}
			/>

			<UploadButton
				value={doc}
				onChange={setDocument}
				instrumentId={instrumentId}
				group='act'
				sx={{
					width: 200,
					borderTopLeftRadius: 0,
					borderBottomLeftRadius: 0,
					borderLeft: 0,
					borderColor: '#c4c4c4',
				}}
			/>
		</Stack>
	)
}
