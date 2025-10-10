import { Button, Divider, Stack, useTheme } from '@mui/material'

import { PermRules } from '@/constants/permissions'
import { useAppDispatch } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { CreateDialog } from '../Dialogs/CreateDialog'
import { ActiveSection } from '@/features/sections/components/Active/Active'
import { ToolsMenu } from '../../modules/tools/components/ToolsMenuLazy'
import { Search } from '../Search/Search'
import { Filters } from '../Filters/Filters'
import { FastSelect } from '../Select/FastSelect'
import { Status } from '../Status/Status'
import { Setting } from '../Setting/SettingLazy'

export const Header = () => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const createHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: true }))
	}

	return (
		<Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'} mt={1} mb={0.5} mx={2}>
			<Stack direction={'row'} sx={{ mr: 1 }}>
				<ActiveSection />

				{useCheckPermission(PermRules.SI.Write) ? (
					<>
						<Status />
						<Divider orientation='vertical' flexItem />

						<Button onClick={createHandler} sx={{ textTransform: 'inherit', ml: 0.5, p: '6px 12px' }}>
							<PlusIcon fontSize={12} mr={1} fill={palette.primary.main} /> Добавить
						</Button>
					</>
				) : null}
			</Stack>

			<Search />

			<Stack direction={'row'} spacing={{ xl: 2, lg: 1 }} sx={{ ml: 1 }}>
				<Setting />
				<FastSelect />
				<Filters />

				<ToolsMenu />
			</Stack>

			<CreateDialog />
		</Stack>
	)
}
