import { IconButton } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { History } from './History'

type Context = string

export const HistoryDialog = () => {
	const modal = useAppSelector(getDialogState('History'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'History', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Посмотреть историю'}
			headerActions={
				<IconButton onClick={closeHandler} size='large' sx={{ fill: '#505050', mr: 2 }}>
					<TimesIcon fontSize={12} />
				</IconButton>
			}
			body={<History instrumentId={context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='lg'
			fullWidth
		/>
	)
}
