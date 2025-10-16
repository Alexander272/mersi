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
import { HistoryDialog } from '../../history/components/Dialog'
import { DeleteLocationDialog } from '../../locations/components/Dialogs/Delete'
import { ForcedReceiptDialog } from '../../locations/components/Dialogs/Forced'
import { CreateOnBase } from './CreateOnBase'
import {
	ChangePosition,
	Edit,
	History,
	Preservation,
	Repair,
	SendToVerification,
	TransferToDep,
	TransferToSave,
	Verification,
	WriteOff,
} from './ContextMenuItems'
import { Cancel, Forced, Location, Receive, Reserve } from './LocationItems'

type ItemProps = { onClick?: () => void; label?: string }
type MenuItem = { el: FC<ItemProps>; action?: DialogVariants }

const ContextMenuItems = new Map<string, MenuItem>([
	['create', { el: CreateOnBase }],
	['edit', { el: Edit, action: 'EditTableItem' }],
	['change-position', { el: ChangePosition, action: 'ChangePosition' }],
	['verification', { el: Verification, action: 'NewVerification' }],
	['sent-to-verification', { el: SendToVerification, action: 'SendToVerification' }],
	['location', { el: Location, action: 'NewLocation' }],
	['reserve', { el: Reserve, action: 'SendToReserve' }],
	['receive', { el: Receive, action: 'Receive' }],
	['forced', { el: Forced, action: 'Forced' }],
	['cancel', { el: Cancel, action: 'DeleteLocation' }],
	['repair-info', { el: Repair, action: 'AddRepair' }],
	['preservation-info', { el: Preservation, action: 'AddPreservation' }],
	['transfer-to-save', { el: TransferToSave, action: 'AddTransferToSave' }],
	['transfer-to-department', { el: TransferToDep, action: 'AddTransferToDep' }],
	['write-off', { el: WriteOff, action: 'WriteOff' }],
	['history', { el: History, action: 'History' }],
])

export const ContextMenu = () => {
	const section = useAppSelector(getSection)
	const contextMenu = useAppSelector(getContextMenu)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetContextMenuQuery({ section: section?.id || '' }, { skip: !section?.id })

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
			<HistoryDialog />
			<DeleteLocationDialog />
			<ForcedReceiptDialog />
		</>
	)
}

export default ContextMenu
