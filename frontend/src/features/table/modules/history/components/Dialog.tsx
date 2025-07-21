import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
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
			body={<History instrumentId={context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='lg'
			fullWidth
		/>
	)
}
