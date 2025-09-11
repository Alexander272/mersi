import { FC } from 'react'
import { ListItemIcon, MenuItem } from '@mui/material'
import dayjs from 'dayjs'

import { PermRules } from '@/constants/permissions'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetLastLocationQuery } from '../../locations/locationsApiSlice'
import { getContextMenu } from '@/features/table/tableSlice'
import { FileSyncIcon } from '@/components/Icons/FileSyncIcon'
import { CancelIcon } from '@/components/Icons/CancelIcon'
import { ExchangeIcon } from '@/components/Icons/ExchangeIcon'

type Props = {
	onClick?: () => void
	label?: string
}

export const Location: FC<Props> = ({ onClick, label }) => {
	const ctx = useAppSelector(getContextMenu)

	const { data } = useGetLastLocationQuery(ctx?.active || '', { skip: !ctx })

	return (
		<MenuItem onClick={onClick} disabled={data?.data.status == 'moved'}>
			<ListItemIcon>
				<ExchangeIcon fontSize={18} fill={'#757575'} />
			</ListItemIcon>
			{label ? label : 'Добавить перемещение'}
		</MenuItem>
	)
}

export const Reserve: FC<Props> = ({ onClick, label }) => {
	const ctx = useAppSelector(getContextMenu)

	const { data } = useGetLastLocationQuery(ctx?.active || '', { skip: !ctx })

	return (
		<MenuItem onClick={onClick} disabled={data?.data.status != 'used'}>
			<ListItemIcon>
				<ExchangeIcon fontSize={18} fill={'#757575'} />
			</ListItemIcon>
			{label ? label : 'Вернуть инструмент'}
		</MenuItem>
	)
}

export const Receive: FC<Props> = ({ onClick, label }) => {
	const ctx = useAppSelector(getContextMenu)

	const { data } = useGetLastLocationQuery(ctx?.active || '', { skip: !ctx })

	const reserve = useCheckPermission(PermRules.Reserve.Write)
	const loc = useCheckPermission(PermRules.Location.Write)
	if ((loc && data?.data.lastPlaceId == '') || (reserve && data?.data.lastPlaceId != '')) return null

	if (data?.data.status != 'moved') return null
	return (
		<MenuItem onClick={onClick}>
			<ListItemIcon>
				<FileSyncIcon fontSize={18} fill={'#757575'} />
			</ListItemIcon>
			{label ? label : 'Получить инструмент'}
		</MenuItem>
	)
}

export const Forced: FC<Props> = ({ onClick, label }) => {
	const ctx = useAppSelector(getContextMenu)

	const { data } = useGetLastLocationQuery(ctx?.active || '', { skip: !ctx })

	if (data?.data.status != 'moved' || data?.data.lastPlaceId != '') return null
	return (
		<MenuItem onClick={onClick}>
			<ListItemIcon>
				<FileSyncIcon fontSize={18} fill={'#757575'} />
			</ListItemIcon>
			{label ? label : 'Отметить инструмент как полученный'}
		</MenuItem>
	)
}

export const Cancel: FC<Props> = ({ onClick, label }) => {
	const ctx = useAppSelector(getContextMenu)

	const { data } = useGetLastLocationQuery(ctx?.active || '', { skip: !ctx })

	const reserve = useCheckPermission(PermRules.Reserve.Write)
	const loc = useCheckPermission(PermRules.Location.Write)
	if ((loc && data?.data.lastPlaceId != '') || (reserve && data?.data.lastPlaceId == '')) return null

	// if (data?.data.lastPlaceId == '') return null
	if (!data?.data || (data.data.dateOfIssue != dayjs().startOf('d').toISOString() && data.data.status != 'moved'))
		return null

	return (
		<MenuItem onClick={onClick}>
			<ListItemIcon>
				<CancelIcon fontSize={18} fill={'#757575'} />
			</ListItemIcon>
			{label ? label : 'Отменить перемещение'}
		</MenuItem>
	)
}
