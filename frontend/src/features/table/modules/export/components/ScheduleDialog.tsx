import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { PeriodForm } from './PeriodForm'

export const ScheduleDialog = () => {
	const modal = useAppSelector(getDialogState('Schedule'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Schedule', isOpen: false }))
	}

	return (
		<Dialog
			title={'Укажите период'}
			body={<PeriodForm />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='sm'
			fullWidth
		/>
	)
}
