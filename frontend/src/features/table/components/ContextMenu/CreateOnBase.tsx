import { CircularProgress, ListItemIcon, MenuItem } from '@mui/material'

import { localKeys } from '../../constants/storage'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useLazyGetInstrumentByIdQuery } from '../../instrumentApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { getContextMenu, setContextMenu } from '../../tableSlice'
import { CopyIcon } from '@/components/Icons/CopyIcon'

export const CreateOnBase = () => {
	const context = useAppSelector(getContextMenu)
	const dispatch = useAppDispatch()

	const [get, { isLoading }] = useLazyGetInstrumentByIdQuery()

	const contextHandler = async () => {
		if (!context?.active) return

		try {
			const data = await get(context?.active).unwrap()

			const form = { instrument: data.data }
			localStorage.setItem(localKeys.form, JSON.stringify(form))

			dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: true }))
			dispatch(setContextMenu())
		} catch {
			/* empty */
		}
	}

	return (
		<MenuItem onClick={contextHandler}>
			<ListItemIcon>
				{isLoading ? <CircularProgress size={18} /> : <CopyIcon fontSize={18} fill={'#757575'} />}
			</ListItemIcon>
			Создать на основании
		</MenuItem>
	)
}
