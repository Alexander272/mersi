import { FC } from 'react'
import { Menu } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetContextMenuQuery } from '../contextMenuApiSlice'
import { changeDialogIsOpen, type DialogVariants } from '@/features/dialog/dialogSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getContextMenu, setContextMenu } from '../../../tableSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { ChangePositionDialog } from '../../../components/Dialogs/ChangePositionDialog'
import { EditDialog } from '../../../components/Dialogs/EditDialog'
import { CreateOnBase } from './CreateOnBase'
import {
	ChangePosition,
	Edit,
	History,
	Preservation,
	Repair,
	TransferToDep,
	TransferToSave,
	Verification,
} from './ContextMenuItems'

type ItemProps = { onClick?: () => void; label?: string }
type MenuItem = { el: FC<ItemProps>; action?: DialogVariants }

const ContextMenuItems = new Map<string, MenuItem>([
	['create', { el: CreateOnBase }],
	['edit', { el: Edit, action: 'EditTableItem' }],
	['change-position', { el: ChangePosition, action: 'ChangePosition' }],
	['verification', { el: Verification, action: 'NewVerification' }],
	['repair-info', { el: Repair, action: 'AddRepair' }],
	['preservation-info', { el: Preservation, action: 'AddPreservation' }],
	['transfer-to-save', { el: TransferToSave, action: 'AddTransferToSave' }],
	['transfer-to-department', { el: TransferToDep, action: 'AddTransferToDep' }],
	['history', { el: History, action: 'History' }],
])

//TODO я не знаю стоит ли добавлять все возможные пункты в контекстное меню, поэтому можно сделать
// в меню инструменты поле для добавления пункта в контекстное меню (что-то типа избранного, сделать звездочку с краю поля)

export const ContextMenu = () => {
	const section = useAppSelector(getSection)
	const contextMenu = useAppSelector(getContextMenu)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetContextMenuQuery(section?.id || '', { skip: !section?.id })

	const closeHandler = () => {
		dispatch(setContextMenu())
	}

	const contextHandler = (variant: DialogVariants) => () => {
		dispatch(changeDialogIsOpen({ variant, isOpen: true, context: contextMenu?.active || '' }))
		closeHandler()
	}

	// const SiMenuItemsWriter = [
	// 	<CreateOnBase key='create' />,
	// 	<MenuItem key='edit' onClick={contextHandler('EditTableItem')}>
	// 		<ListItemIcon>
	// 			<EditIcon fontSize={16} fill={'#757575'} />
	// 		</ListItemIcon>
	// 		Редактировать
	// 	</MenuItem>,
	// 	<MenuItem key={'change-position'} onClick={contextHandler('ChangePosition')}>
	// 		<ListItemIcon>
	// 			<TransferIcon fontSize={20} fill={'#757575'} />
	// 		</ListItemIcon>
	// 		Изменить номер позиции
	// 	</MenuItem>,
	// ]

	return (
		<>
			{isFetching && <BoxFallback />}

			<Menu
				open={Boolean(contextMenu)}
				onClose={closeHandler}
				anchorReference='anchorPosition'
				anchorPosition={
					contextMenu ? { top: contextMenu.coords.mouseY, left: contextMenu.coords.mouseX } : undefined
				}
			>
				{data?.data.map(d => {
					const item = ContextMenuItems.get(d.name)
					const Elem = item?.el
					if (!Elem) return null
					return (
						<Elem
							key={d.id}
							onClick={item.action ? contextHandler(item.action) : undefined}
							label={d.label}
						/>
					)
				})}
			</Menu>

			<EditDialog />
			<ChangePositionDialog />
		</>
	)
}

export default ContextMenu
