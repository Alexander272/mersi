import { EditIcon } from '@/components/Icons/EditIcon'
import { TransferIcon } from '@/components/Icons/TransferIcon'
import { ListItemIcon, MenuItem } from '@mui/material'
import { FC } from 'react'

type Props = {
	onClick?: () => void
	label?: string
}

export const Edit: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<EditIcon fontSize={16} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Редактировать'}
	</MenuItem>
)

export const ChangePosition: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<TransferIcon fontSize={20} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Изменить номер позиции'}
	</MenuItem>
)
