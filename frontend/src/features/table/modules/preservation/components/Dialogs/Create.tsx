import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Forms/Create'

type Context = string | string[]

export const CreatePreservationDialog = () => {
	const modal = useAppSelector(getDialogState('AddPreservation'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'AddPreservation', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Добавить сведения о консервации'}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
