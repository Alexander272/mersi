import { ChangeEvent, FC } from 'react'
import { Box, Button, CircularProgress, SxProps, Theme, Tooltip } from '@mui/material'
import { toast } from 'react-toastify'

import type { IDocument } from '../../types/document'
import { AcceptedFiles } from '@/constants/accept'
import { UploadIcon } from '@/components/Icons/UploadIcon'
import { QuestionIcon } from '@/components/Icons/QuestionIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { useDeleteFileMutation, useUploadFilesMutation } from '../../fileApiSlice'
import { FileIcon } from '../FileIcon/FileIcon'
import Input from './Input'

type Props = {
	value: IDocument | null
	onChange: (value: IDocument | null) => void
	instrumentId: string
	group: string
	sx?: SxProps<Theme>
}

export const UploadButton: FC<Props> = ({ value, onChange, instrumentId, group, sx }) => {
	const [upload, { isLoading: uploading }] = useUploadFilesMutation()
	const [remove, { isLoading: removing }] = useDeleteFileMutation()

	const changeHandler = (event: ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files
		if (!files) return
		const acceptedFiles: File[] = []

		for (let i = 0; i < files.length; i++) {
			const file = files[i]

			if (!(file.type in AcceptedFiles)) {
				toast.error(`Файл ${file.name} имеет неразрешенный тип`)
				continue
			}
			acceptedFiles.push(file)
		}
		uploadFiles(acceptedFiles)
	}

	const uploadFiles = async (files: File[]) => {
		if (!files) return
		const data = new FormData()
		data.append('instrumentId', instrumentId)
		data.append('group', group)
		files.forEach((file: File) => data.append('files', file))

		try {
			const res = await upload({ data }).unwrap()
			onChange(res.data[0])
		} catch {
			/* empty */
		}
	}

	const removeFiles = async () => {
		if (!value) return
		const data = {
			instrumentId,
			group,
			id: value.id,
			filename: value.label,
		}
		try {
			await remove(data).unwrap()
			onChange(null)
		} catch {
			// empty
		}
	}

	return (
		<Button variant='outlined' color='inherit' sx={{ position: 'relative', ...sx }} component='label'>
			{uploading || removing ? (
				<CircularProgress size={24} />
			) : value ? (
				<FileIcon type={value.type} />
			) : (
				<UploadIcon />
			)}

			<Input onChange={changeHandler} type='file' disabled={Boolean(value)} />

			{value ? (
				<Box
					onClick={removeFiles}
					position={'absolute'}
					right={8}
					top={'50%'}
					padding={1}
					height={28}
					borderRadius={8}
					display={'flex'}
					justifyContent={'center'}
					alignItems={'center'}
					sx={{
						cursor: 'pointer',
						transition: '.3s all ease-in-out',
						transform: 'translateY(-50%)',
						':hover': { backgroundColor: '#c5c5c5' },
					}}
				>
					<TimesIcon fontSize={12} />
				</Box>
			) : (
				<Tooltip
					title='Допустимые форматы: .doc, .docx, .odt, .xls, .xlsx, .pdf, .png, .jpeg, .jpg, .csv'
					enterDelay={100}
					arrow
				>
					<Box
						position={'absolute'}
						right={8}
						top={6}
						padding={0.5}
						height={26}
						borderRadius={2}
						sx={{
							cursor: 'help',
							transition: '.3s all ease-in-out',
							':hover': { backgroundColor: '#dfdfdf' },
						}}
					>
						<QuestionIcon fontSize={16} color='#828282' />
					</Box>
				</Tooltip>
			)}
		</Button>
	)
}
