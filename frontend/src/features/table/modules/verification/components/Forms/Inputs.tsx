import { FC, useEffect, useState } from 'react'
import { FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFieldArray, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IDocument } from '@/features/files/types/document'
import { VerificationStatuses } from '../../constants/status'
import { useAppSelector } from '@/hooks/redux'
import { useGetTempFilesQuery } from '@/features/files/fileApiSlice'
import { useGetVerificationFieldsQuery } from '../../modules/fieldsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { UploadButton } from '@/features/files/components/UploadButton/UploadButton'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Upload } from '@/features/files/components/Upload/Upload'

const docsGroup = 'verifications'

type Props = {
	instrumentId: string
	interval?: number
	minDate?: string
}

export const Inputs: FC<Props> = ({ instrumentId, interval, minDate = '2000-01-01' }) => {
	const section = useAppSelector(getSection)

	const { data, isFetching } = useGetVerificationFieldsQuery(
		{ section: section?.id || '', group: 'form' },
		{ skip: !section?.id }
	)

	return (
		<Stack spacing={2} mb={2}>
			{isFetching && <BoxFallback />}

			{data?.data.map(f => {
				if (f.type == 'date')
					return (
						<DateField key={f.id} field={f.field} label={f.label} minDate={minDate} interval={interval} />
					)
				if (f.field == 'status') return <StatusField key={f.id} field={f.field} label={f.label} />
				if (f.field == 'registerLink') return <LinkField key={f.id} field={f.field} label={f.label} />
				if (f.field == 'notes') return <NotesField key={f.id} field={f.field} label={f.label} />
				if (f.type == 'file')
					return <FileField key={f.id} field={f.field} label={f.label} instrumentId={instrumentId} />
				if (f.type == 'files')
					return <FilesBoxField key={f.id} field={f.field} label={f.label} instrumentId={instrumentId} />
			})}
		</Stack>
	)
}

type FieldProps = {
	label: string
	field: string
	instrumentId?: string
}

const DateField: FC<FieldProps & { minDate: string; interval?: number }> = ({ label, field, minDate, interval }) => {
	const { control, watch, setValue } = useFormContext()

	let date = ''
	let status = ''
	if (field == 'nextVerificationDate') {
		date = watch('verificationDate')
		status = watch('status')
	}

	useEffect(() => {
		if (field == 'nextVerificationDate') {
			console.log(date, interval)
			if (interval && date) {
				const newDate = dayjs(date).add(interval, 'M').subtract(1, 'd').toISOString()
				console.log(newDate)

				setValue('nextVerificationDate', newDate)
			}
		}
	}, [date, field, interval, setValue])

	if (field == 'nextVerificationDate' && status != VerificationStatuses.Work) return null
	return (
		<Controller
			control={control}
			name={field}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={field.value ? dayjs(field.value) : null}
					onChange={value => field.onChange(value?.startOf('d').toISOString())}
					label={label}
					showDaysOutsideCurrentMonth
					fixedWeekNumber={6}
					minDate={dayjs(minDate)}
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

	const { data, isFetching } = useGetTempFilesQuery(
		{ group: docsGroup, instrument: instrumentId },
		{ skip: !instrumentId }
	)

	useEffect(() => {
		if (data?.data.length) {
			setDoc(data.data[0])
			setValue('doc', data.data[0]?.label || '')
			setValue(`docId`, data.data[0]?.id || '')
		}
	}, [data, setValue])

	const setDocument = (value: IDocument | null) => {
		setDoc(value)
		setValue(field, value?.label || '')
		setValue(`${field}Id`, value?.id || '')
	}

	return (
		<Stack direction={'row'}>
			{isFetching && <BoxFallback />}

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
				group={docsGroup}
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

const FilesBoxField: FC<FieldProps> = ({ field, instrumentId = '' }) => {
	const [doc, setDoc] = useState<IDocument[]>([])

	const { control } = useFormContext()
	const { replace } = useFieldArray({ control, name: field as 'docs' })

	const { data, isFetching } = useGetTempFilesQuery(
		{ group: docsGroup, instrument: instrumentId },
		{ skip: !instrumentId }
	)

	useEffect(() => {
		if (data?.data.length) {
			setDoc(data.data)
			replace(data.data.map(d => ({ docId: d.id, doc: d.label })))
		}
	}, [data, field, replace])

	const setDocument = (value: IDocument[]) => {
		setDoc(value)
		replace(value.map(d => ({ docId: d.id, doc: d.label })) || [])
	}

	if (isFetching) return <BoxFallback />
	return <Upload value={doc} onChange={setDocument} instrumentId={instrumentId} group={docsGroup} />
}
