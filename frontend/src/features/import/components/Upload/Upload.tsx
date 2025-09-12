import { type ChangeEvent, type DragEvent, type FC, useState } from 'react'
import { Box, type SxProps, type Theme, Typography } from '@mui/material'

import { ImageIcon } from './ImageIcon'
import Input from './styled/Input'

type Props = {
	fullscreen?: boolean
	sx?: SxProps<Theme>
	value?: string
	onChange: (file?: File | null) => void
}

export const Upload: FC<Props> = props => {
	return <UploadInput {...props} />
}

const UploadInput: FC<Props> = ({ onChange, sx }) => {
	const [hasDropFiles, setHasDropFiles] = useState(false)

	const changeHandler = (event: ChangeEvent<HTMLInputElement>) => {
		const files = event.target.files
		if (!files) return

		onChange(files[0])
	}

	const dragHandler = (event: DragEvent<HTMLDivElement>) => {
		event.preventDefault()
		event.stopPropagation()

		if (event.type === 'dragenter' || event.type === 'dragover') {
			setHasDropFiles(true)
		} else if (event.type === 'dragleave') {
			setHasDropFiles(false)
		}
	}

	const dropHandler = (event: DragEvent<HTMLDivElement>) => {
		event.preventDefault()
		event.stopPropagation()
		setHasDropFiles(false)

		const files = event.dataTransfer.files
		onChange(files[0])
	}

	return (
		<Box
			// component='label'

			display={'flex'}
			flexDirection={'column'}
			justifyContent={'center'}
			alignItems={'center'}
			flexGrow={1}
			boxShadow={hasDropFiles ? 'inset 0 0 20px #00000028' : undefined}
			borderRadius={3}
			onDragEnter={dragHandler}
			onDragLeave={dragHandler}
			onDragOver={dragHandler}
			onDrop={dropHandler}
		>
			<Box
				component='label'
				border={`1px solid #0000003b`}
				borderRadius={3}
				display={'flex'}
				flexDirection={'column'}
				justifyContent={'center'}
				alignItems={'center'}
				flexGrow={1}
				padding={2}
				sx={{
					cursor: 'pointer',
					transition: 'all 0.3s ease-in-out',
					':hover': {
						// boxShadow: 'inset 0 0 20px #00000028',
						borderColor: '#000000de',
					},
					...sx,
				}}
			>
				<ImageIcon fontSize={100} fill={'#646464'} />
				<Typography fontSize={'1.1rem'} color={'#464646'}>
					Выберите или перетащите файл
				</Typography>
				<Input onChange={changeHandler} type='file' />
			</Box>
		</Box>
	)
}
