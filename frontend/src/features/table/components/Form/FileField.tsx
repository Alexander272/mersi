import { FC, useState } from 'react'
import { Stack, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import type { IDocument } from '@/features/files/types/document'
import { UploadButton } from '@/features/files/components/UploadButton/UploadButton'

type Props = {
	data: ICreateFormField
	instrumentId?: string
}

export const FileField: FC<Props> = ({ data, instrumentId = '' }) => {
	const [doc, setDoc] = useState<IDocument | null>(null)

	const { control, setValue, watch } = useFormContext()
	const watchField = watch(data.hide)

	const setDocument = (value: IDocument | null) => {
		setDoc(value)
		setValue(`${data.path}.${data.field}`, value?.label || '')
		setValue(`${data.path}.${data.field}Id`, value?.id || '')
	}

	if (data.hide && watchField) return null
	return (
		<Stack direction={'row'}>
			<Controller
				control={control}
				name={data.path + '.' + data.field}
				rules={{ required: data.isRequired }}
				render={({ field, fieldState: { error } }) => (
					<TextField
						{...field}
						value={field.value || ''}
						label={data.fieldName}
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
					//TODO можно еще попробовать сделать границу прозрачной и добавить отрицательный отступ, чтобы выделять границу при наведении
					borderLeft: 0,
					borderColor: '#c4c4c4',
				}}
			/>
		</Stack>
	)
}
