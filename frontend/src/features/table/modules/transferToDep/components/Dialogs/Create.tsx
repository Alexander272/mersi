import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Forms/Create'

type Context = string | string[]

export const CreateTransferToDepartmentDialog = () => {
	const modal = useAppSelector(getDialogState('AddTransferToDep'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'AddTransferToDep', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Добавить сведения о передаче в другое подразделение'}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
