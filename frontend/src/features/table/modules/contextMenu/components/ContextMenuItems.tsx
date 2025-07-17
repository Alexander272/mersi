import { FC } from 'react'
import { ListItemIcon, MenuItem } from '@mui/material'

import { EditIcon } from '@/components/Icons/EditIcon'
import { TransferIcon } from '@/components/Icons/TransferIcon'
import { VerifyIcon } from '@/components/Icons/VerifyIcon'
import { RepairIcon } from '@/components/Icons/RepairIcon'
import { ToolboxIcon } from '@/components/Icons/ToolboxIcon'
import { ProductReplace } from '@/components/Icons/ProductReplace'
import { ProductReturn } from '@/components/Icons/ProductReturn'
import { HistoryIcon } from '@/components/Icons/HistoryIcon'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'

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

export const Verification: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<VerifyIcon fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Добавить поверку'}
	</MenuItem>
)

export const Repair: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<RepairIcon fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Добавить сведения о ремонте'}
	</MenuItem>
)

export const Preservation: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<ToolboxIcon fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Добавить сведения о консервации'}
	</MenuItem>
)

export const TransferToSave: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<ProductReplace fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Добавить сведения о передаче на хранение'}
	</MenuItem>
)

export const TransferToDep: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<ProductReturn fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Добавить сведения о передаче другому подразделению'}
	</MenuItem>
)

export const WriteOff: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<FileDeleteIcon fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'Списать'}
	</MenuItem>
)

export const History: FC<Props> = ({ onClick, label }) => (
	<MenuItem onClick={onClick}>
		<ListItemIcon>
			<HistoryIcon fontSize={18} fill={'#757575'} />
		</ListItemIcon>
		{label ? label : 'История'}
	</MenuItem>
)
