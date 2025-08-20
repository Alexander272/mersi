import { ChangeEvent, FC, MouseEvent, useEffect, useState } from 'react'
import {
	Box,
	CircularProgress,
	IconButton,
	List,
	ListItem,
	ListItemIcon,
	ListItemSecondaryAction,
	ListItemText,
	Stack,
	Tooltip,
	Typography,
	useTheme,
} from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { convertFileSize } from '@/features/files/utils/convertFileSize'
import { useDeleteFileMutation, useGetFileListQuery, useUploadFilesMutation } from '@/features/files/fileApiSlice'
import { AcceptedFiles } from '@/constants/accept'
import { UploadIcon } from '@/components/Icons/UploadIcon'
import { QuestionIcon } from '@/components/Icons/QuestionIcon'
import { PdfIcon } from '@/components/Icons/PdfIcon'
import { DocIcon } from '@/components/Icons/DocIcon'
import { ImageIcon } from '@/components/Icons/ImageIcon'
import { SheetIcon } from '@/components/Icons/SheetIcon'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import Input from './Input'
import Button from './Button'

type Props = {
	instrumentId: string
	verificationId: string
}

const Types = {
	doc: <DocIcon />,
	pdf: <PdfIcon />,
	image: <ImageIcon />,
	sheet: <SheetIcon />,
}

export const Upload: FC<Props> = ({ verificationId, instrumentId }) => {
	const [files, setFiles] = useState<File[]>([])

	const { palette } = useTheme()

	const { data } = useGetFileListQuery(
		{ verificationId: verificationId || '', instrumentId: instrumentId }
		// { skip: !instrumentId }
	)

	const [upload, { isSuccess, isError }] = useUploadFilesMutation()
	const [remove] = useDeleteFileMutation()

	useEffect(() => {
		if (isSuccess || isError) setFiles([])
	}, [isError, isSuccess])

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
		setFiles(acceptedFiles)
		void uploadFiles(acceptedFiles)
	}

	const uploadFiles = async (files: File[]) => {
		const data = new FormData()
		files.forEach((file: File) => data.append('files', file))
		data.append('instrumentId', instrumentId)
		data.append('verificationId', verificationId)

		await upload({ data, verificationId }).unwrap()
	}

	const deleteHandler = (event: MouseEvent<HTMLButtonElement>) => {
		const { id } = (event.target as HTMLButtonElement).dataset
		const doc = data?.data.find(d => d.id == id)
		if (!doc) return

		void deleteFile(doc)
	}

	const deleteFile = async (doc: IDocument) => {
		const data = {
			instrumentId,
			verificationId,
			id: doc.id,
			filename: doc.label,
		}
		try {
			await remove(data).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack direction={'row'} spacing={3} alignItems={'flex-start'}>
			<Button component='label'>
				<UploadIcon />
				<Typography ml={1}>Загрузить файлы</Typography>

				<Input onChange={changeHandler} type='file' multiple />
			</Button>

			<Stack
				spacing={1}
				flexGrow={1}
				width={'100%'}
				border={'1px dashed #c4c4c4'}
				borderRadius={3}
				paddingY={0.75}
				paddingX={2}
				minHeight={38}
				position={'relative'}
			>
				<Tooltip
					title='Допустимые форматы: .doc, .docx, .odt, .xls, .xlsx, .pdf, .png, .jpeg, .jpg, .csv'
					enterDelay={100}
					arrow
				>
					<Box
						position={'absolute'}
						right={8}
						top={4}
						padding={0.5}
						height={26}
						borderRadius={2}
						sx={{
							cursor: 'help',
							transition: '.3s all ease-in-out',
							':hover': { backgroundColor: '#eee' },
						}}
					>
						<QuestionIcon fontSize={16} color='#828282' />
					</Box>
				</Tooltip>

				{/* //TODO как-то файлы дергаются */}
				<List dense disablePadding sx={{ mt: '0!important' }}>
					{(data?.data || []).map(d => (
						<ListItem key={d.id} sx={{ paddingY: 0, pl: 1 }}>
							<ListItemIcon sx={{ minWidth: 40 }}>{Types[d.type as 'doc']}</ListItemIcon>
							<ListItemText primary={d.label} secondary={convertFileSize(d.size, 2) + ' МБ'} />
							<ListItemSecondaryAction>
								<IconButton data-id={d.id} onClick={deleteHandler}>
									<DeleteIcon fontSize={18} color={palette.error.main} pointerEvents='none' />
								</IconButton>
							</ListItemSecondaryAction>
						</ListItem>
					))}
					{files.map(d => (
						<ListItem key={d.name + d.size} sx={{ paddingY: 0, pl: 1 }}>
							<ListItemIcon sx={{ minWidth: 40 }}>
								<CircularProgress size={24} />
							</ListItemIcon>
							<ListItemText primary={d.name} secondary={convertFileSize(d.size, 2) + ' МБ'} />
						</ListItem>
					))}
				</List>
			</Stack>
		</Stack>
	)
}
