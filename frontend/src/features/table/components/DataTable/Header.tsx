import { Button, Stack, useTheme } from '@mui/material'

import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { CreateDialog } from '../Dialogs/CreateDialog'
import { ActiveSection } from '@/features/sections/components/Active/Active'
import { ToolsMenu } from '../../modules/tools/components/ToolsMenuLazy'
import { Search } from '../Search/Search'
import { Filters } from '../Filters/Filters'
import { FastSelect } from '../Select/FastSelect'
import { Status } from '../Status/Status'

export const Header = () => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const createHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: true }))
	}

	return (
		<Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'} mt={1} mb={0.5} mx={2}>
			<Stack direction={'row'} spacing={1}>
				<ActiveSection />

				{/* active list */}
				<Status />

				<Button onClick={createHandler} variant='outlined'>
					<PlusIcon fontSize={12} mr={1} fill={palette.primary.main} /> Добавить
				</Button>
			</Stack>

			<Search />

			<Stack direction={'row'} spacing={2}>
				{/* <Setting />*/}
				<FastSelect />
				<Filters />

				<ToolsMenu />
			</Stack>

			<CreateDialog />
		</Stack>
	)
}
