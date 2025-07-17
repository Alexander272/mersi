import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Forms/Create'

type Context = string | string[]

export const CreateWriteOffDialog = () => {
	const modal = useAppSelector(getDialogState('WriteOff'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'WriteOff', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Добавить сведения о списании'}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
