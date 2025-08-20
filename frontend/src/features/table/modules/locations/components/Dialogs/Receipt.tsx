import { IconButton } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { ReceiptOne } from '../Forms/Receipt/ReceiptOne'
import { ReceiptMany } from '../Forms/Receipt/ReceiptMany'

type Context = string | string[]

export const ReceiptDialog = () => {
	const modal = useAppSelector(getDialogState('Receive'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Receive', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Подтвердить получение инструментов'}
			headerActions={
				<IconButton onClick={closeHandler} size='large' sx={{ fill: '#505050', mr: 2 }}>
					<TimesIcon fontSize={12} />
				</IconButton>
			}
			body={typeof context == 'string' ? <ReceiptOne id={context} /> : <ReceiptMany />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth={typeof context == 'string' ? 'sm' : 'md'}
			fullWidth
		/>
	)
}
