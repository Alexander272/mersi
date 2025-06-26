import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { EditForm } from './EditForm'

type Context = string

export const EditDialog = () => {
	const modal = useAppSelector(getDialogState('EditTableItem'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditTableItem', isOpen: false }))
	}

	return (
		<Dialog
			title={'Редактировать'}
			body={<EditForm id={(modal?.context as Context) || ''} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
