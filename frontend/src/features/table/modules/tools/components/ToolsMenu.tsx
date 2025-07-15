import { useRef, useState } from 'react'
import { Button, Menu } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetToolsMenuQuery } from '../toolsMenuApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { changeDialogIsOpen, DialogVariants } from '@/features/dialog/dialogSlice'
import { getSelected } from '@/features/table/tableSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { CreateRepairDialog } from '../../repair/components/Dialogs/Create'
import { CreateVerificationDialog } from '../../verification/components/Dialogs/Create'
import { MenuItems } from './ToolsMenuItem'
import { CreatePreservationDialog } from '../../preservation/components/Dialogs/Create'
import { CreateTransferToSaveDialog } from '../../transferToSave/components/Dialogs/Create'

export const ToolsMenu = () => {
	const anchor = useRef<HTMLButtonElement>(null)
	const [open, setOpen] = useState(false)

	const selected = useAppSelector(getSelected)
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetToolsMenuQuery(section?.id || '', { skip: !section?.id })

	const toggleHandler = () => setOpen(prev => !prev)

	const menuHandler = (variant: DialogVariants) => () => {
		dispatch(changeDialogIsOpen({ variant, isOpen: true, context: Object.keys(selected) }))
		toggleHandler()
	}

	return (
		<>
			<Button
				ref={anchor}
				onClick={toggleHandler}
				variant='outlined'
				color='inherit'
				// size='small'
				// sx={{ paddingX: 1.5 }}
			>
				Инструменты
			</Button>

			<Menu
				open={open}
				onClose={toggleHandler}
				anchorEl={anchor.current}
				transformOrigin={{ horizontal: 'center', vertical: 'top' }}
				anchorOrigin={{ horizontal: 'center', vertical: 'bottom' }}
				slotProps={{
					paper: {
						elevation: 0,
						sx: {
							overflow: 'visible',
							filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
							mt: 1.5,
							maxWidth: 550,
							width: '100%',
							'&:before': {
								content: '""',
								display: 'block',
								position: 'absolute',
								top: 0,
								right: '20%',
								width: 10,
								height: 10,
								bgcolor: 'background.paper',
								transform: 'translate(-65%, -50%) rotate(45deg)',
								zIndex: 0,
							},
						},
					},
				}}
			>
				{isFetching && <BoxFallback />}

				{data?.data.map(d => {
					const item = MenuItems.get(d.name)
					if (item) item.label = d.label
					const Elem = item?.el
					if (!Elem) return null
					return <Elem key={d.id} item={d} onClick={item.action ? menuHandler(item.action) : undefined} />
				})}
			</Menu>

			<CreateRepairDialog />
			<CreateVerificationDialog title={MenuItems.get('verification')?.label || ''} />
			<CreatePreservationDialog />
			<CreateTransferToSaveDialog />
		</>
	)
}
export default ToolsMenu
