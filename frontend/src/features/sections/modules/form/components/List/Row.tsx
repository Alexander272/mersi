import { FC } from 'react'
import { IconButton, Stack, Typography } from '@mui/material'

import type { ICreateFormField } from '../../types/create'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { DragIcon } from '@/components/Icons/DragIcon'
import { EditIcon } from '@/components/Icons/EditIcon'

type Props = {
	id: string
	data: ICreateFormField
}

export const Row: FC<Props> = ({ data, id }) => {
	const dispatch = useAppDispatch()

	const openEditDialog = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: true, context: { item: { ...data, id } } }))
	}

	return (
		<Stack
			direction={'row'}
			alignItems={'center'}
			mb={1}
			position={'relative'}
			sx={{
				':after': {
					content: '""',
					position: 'absolute',
					left: 40,
					bottom: 0,
					height: '1px',
					width: 'calc(100% - 60px)',
					background: '#a0a6b7a3',
				},
			}}
		>
			<IconButton sx={{ cursor: 'grab', mr: 1 }} className='drag'>
				<DragIcon fill={'#a8a8a8'} fontSize={24} />
			</IconButton>

			<Typography>{data.fieldName}</Typography>
			<IconButton onClick={openEditDialog} size='large' sx={{ ml: 0.5 }}>
				<EditIcon fontSize={12} fill={'#6e6e6e'} />
			</IconButton>
		</Stack>
	)
}
