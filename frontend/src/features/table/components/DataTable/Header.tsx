import { Button, Stack, Typography, useTheme } from '@mui/material'

import { PlusIcon } from '@/components/Icons/PlusIcon'
import { useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'

export const Header = () => {
	const { palette } = useTheme()
	const section = useAppSelector(getSection)
	// const dispatch = useAppDispatch()

	const createHandler = () => {
		// dispatch(changeModalIsOpen({ variant: 'create', isOpen: true }))
	}

	return (
		<Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'} mt={1} mb={0.5} mx={2}>
			<Stack direction={'row'}>
				<Typography color={'primary'} variant='h5'>
					{section?.name}
				</Typography>

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
		</Stack>
	)
}
