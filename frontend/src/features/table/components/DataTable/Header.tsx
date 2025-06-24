import { Button, Stack, useTheme } from '@mui/material'

import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { CreateDialog } from '../CreateDialog/CreateDialog'
import { ActiveSection } from '@/features/sections/components/Active/Active'

export const Header = () => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const createHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: true }))
	}

	return (
		<Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'} mt={1} mb={0.5} mx={2}>
			<Stack direction={'row'}>
				<ActiveSection />

				<Button onClick={createHandler} variant='outlined'>
					<PlusIcon fontSize={12} mr={1} fill={palette.primary.main} /> Добавить
				</Button>

				{/* <Search /> */}
			</Stack>

			<Stack direction={'row'} alignItems={'center'} spacing={2}>
				{/* <Setting />
				<Filters /> */}

				{/* <Create /> */}
			</Stack>

			<CreateDialog />
		</Stack>
	)
}
