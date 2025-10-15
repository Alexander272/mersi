import { FC } from 'react'
import { IconButton, Stack, Typography } from '@mui/material'

import type { IStatusForm } from '../../types/status'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { DragIcon } from '@/components/Icons/DragIcon'
import { EditIcon } from '@/components/Icons/EditIcon'

type Props = {
	data: IStatusForm
}

export const StatusItem: FC<Props> = ({ data }) => {
	const dispatch = useAppDispatch()

	const openEditDialog = () => {
		const context = { data }
		dispatch(changeDialogIsOpen({ variant: 'EditStatus', isOpen: true, context }))
	}

	let background = 'inherit'
	if (data.status === 'deleted') background = '#ff8282'
	if (data.status === 'updated') background = '#ffff9b'
	if (data.status === 'new') background = '#a2ffa2'
	if (data.status === 'moved') background = '#e6e6ff'

	return (
		<Stack
			direction={'row'}
			alignItems={'center'}
			mb={1}
			position={'relative'}
			sx={{
				backgroundColor: background,
				borderRadius: 3,
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

			{/* <Typography sx={{ width: '50%', maxWidth: 300 }}>{data.value}</Typography> */}
			<Typography sx={{ ml: 2, textDecoration: data.status === 'deleted' ? 'line-through' : 'none' }}>
				{data.label} ({data.value})
			</Typography>

			<IconButton onClick={openEditDialog} size='large' sx={{ ml: 'auto', mr: 2 }}>
				<EditIcon fontSize={12} fill={'#6e6e6e'} />
			</IconButton>
		</Stack>
	)
}
