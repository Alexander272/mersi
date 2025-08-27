import { FC, useEffect, useState } from 'react'
import { Stack, TextField } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { IDocument } from '@/features/files/types/document'
import { useGetTempFilesQuery } from '@/features/files/fileApiSlice'
import { UploadButton } from '@/features/files/components/UploadButton/UploadButton'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

const min = 1262286000
const docsGroup = 'writeOff'

type Props = {
	instrumentId: string
}

export const Inputs: FC<Props> = ({ instrumentId }) => {
	const [doc, setDoc] = useState<IDocument | null>(null)

	const { control, setValue } = useFormContext()

	const { data, isFetching } = useGetTempFilesQuery(
		{ group: docsGroup, instrument: instrumentId },
		{ skip: !instrumentId }
	)

	useEffect(() => {
		if (data?.data.length) {
			setDoc(data.data[0])
			setValue('docName', data.data[0]?.label || '')
			setValue(`docId`, data.data[0]?.id || '')
		}
	}, [data, setValue])

	const setDocument = (value: IDocument | null) => {
		setDoc(value)
		setValue('docName', value?.label || '')
		setValue(`docId`, value?.id || '')
	}

	return (
		<Stack spacing={2} mb={2}>
			{isFetching && <BoxFallback />}

			<Controller
				control={control}
				name={'date'}
				render={({ field, fieldState: { error } }) => (
					<DatePicker
						{...field}
						value={dayjs(field.value * 1000)}
						onChange={value => field.onChange(value?.startOf('d').unix())}
						label={`Дата передачи`}
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

			<Stack direction={'row'}>
				<Controller
					control={control}
					name={'docName'}
					render={({ field, fieldState: { error } }) => (
						<TextField
							{...field}
							value={field.value || ''}
							label={'Акт о передаче'}
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

			<Controller
				control={control}
				name={'notes'}
				render={({ field }) => <TextField {...field} label={'Примечание'} multiline minRows={4} />}
			/>
		</Stack>
	)
}
