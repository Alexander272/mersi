import { ListItemIcon, Menu, MenuItem } from '@mui/material'

import { PermRules } from '@/constants/permissions'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { changeDialogIsOpen, type DialogVariants } from '@/features/dialog/dialogSlice'
import { getContextMenu, setContextMenu } from '../../tableSlice'
import { EditIcon } from '@/components/Icons/EditIcon'
import { TransferIcon } from '@/components/Icons/TransferIcon'
import { EditDialog } from '../Dialogs/EditDialog'
import { ChangePositionDialog } from '../Dialogs/ChangePositionDialog'
import { CreateOnBase } from './CreateOnBase'

export const ContextMenu = () => {
	const contextMenu = useAppSelector(getContextMenu)
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(setContextMenu())
	}

	const contextHandler = (variant: DialogVariants) => () => {
		dispatch(changeDialogIsOpen({ variant, isOpen: true, context: contextMenu?.active || '' }))
		closeHandler()
	}

	const SiMenuItemsWriter = [
		<CreateOnBase key='create' />,
		<MenuItem key='edit' onClick={contextHandler('EditTableItem')}>
			<ListItemIcon>
				<EditIcon fontSize={16} fill={'#757575'} />
			</ListItemIcon>
			Редактировать
		</MenuItem>,
		<MenuItem key={'change-position'} onClick={contextHandler('ChangePosition')}>
			<ListItemIcon>
				<TransferIcon fontSize={20} fill={'#757575'} />
			</ListItemIcon>
			Изменить номер позиции
		</MenuItem>,
	]

	return (
		<>
			<Menu
				open={Boolean(contextMenu)}
				onClose={closeHandler}
				anchorReference='anchorPosition'
				anchorPosition={
					contextMenu ? { top: contextMenu.coords.mouseY, left: contextMenu.coords.mouseX } : undefined
				}
			>
				{useCheckPermission(PermRules.SI.Write) ? SiMenuItemsWriter : null}
			</Menu>

			<EditDialog />
			<ChangePositionDialog />
		</>
	)
}

export default ContextMenu
