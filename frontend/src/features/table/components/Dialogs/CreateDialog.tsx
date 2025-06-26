import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { CreateForm } from './CreateForm'

type Context = string

export const CreateDialog = () => {
	const modal = useAppSelector(getDialogState('CreateTableItem'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: false }))
	}

	return (
		<Dialog
			title={'Добавить'}
			body={<CreateForm id={(modal?.context as Context) || ''} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
