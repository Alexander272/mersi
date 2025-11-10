import { FC } from 'react'
import { IconButton, useTheme } from '@mui/material'

import type { IHistoryForm } from '../../types/historyTypes'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { HistoryForm } from './Form'

type Context = { data?: IHistoryForm }

type Props = {
	submit: (data: IHistoryForm) => void
}

export const HistoryDialog: FC<Props> = ({ submit }) => {
	const modal = useAppSelector(getDialogState('EditHistory'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditHistory', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.data ? 'Редактировать историю' : 'Добавить историю'}
			headerActions={<DeleteComponent {...(modal?.context as Context)} submit={submit} />}
			body={<HistoryForm {...context} submit={submit} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const DeleteComponent: FC<{ data?: IHistoryForm; submit: (data: IHistoryForm) => void }> = ({ data, submit }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const deleteHandler = async () => {
		if (!data) return
		if (data.id) submit({ ...data, status: 'deleted' })
		else submit({ ...data, status: undefined })
		dispatch(changeDialogIsOpen({ variant: 'EditHistory', isOpen: false }))
	}

	return (
		<Confirm
			buttonComponent={
				<IconButton size='large' sx={{ fill: '#505050', ':hover': { fill: palette.error.main } }}>
					<DeleteIcon fontSize={16} fill={'inherit'} transition={'all 0.2s ease-in-out'} />
				</IconButton>
			}
			width='56px'
			onClick={deleteHandler}
			confirmText={`Вы уверены, что хотите удалить пункт "${data?.label}" в отображении истории?`}
		/>
	)
}
