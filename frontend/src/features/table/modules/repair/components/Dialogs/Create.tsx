import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Form/Create'

type Context = string | string[]

export const CreateRepairDialog = () => {
	const modal = useAppSelector(getDialogState('AddRepair'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'AddRepair', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={'Добавить ремонт'}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
