import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { ChangePositionForm } from './ChangePositionForm'

type Context = string

export const ChangePositionDialog = () => {
	const modal = useAppSelector(getDialogState('ChangePosition'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'ChangePosition', isOpen: false }))
	}

	return (
		<Dialog
			title={'Изменить номер позиции'}
			body={<ChangePositionForm id={(modal?.context as Context) || ''} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
