import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Forms/Create'

type Context = string | string[]

export const CreateTransferToSaveDialog = () => {
	const modal = useAppSelector(getDialogState('AddTransferToSave'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'AddTransferToSave', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Добавить сведения о передаче на сохранение'}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
