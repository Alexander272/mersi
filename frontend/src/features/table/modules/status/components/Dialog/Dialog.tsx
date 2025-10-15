import { FC } from 'react'
import { IconButton, useTheme } from '@mui/material'

import type { IStatusForm } from '../../types/status'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { Form } from '../Form/Form'

type Context = { data?: IStatusForm }

export const StatusDialog: FC<{ submit: (data: IStatusForm) => void }> = ({ submit }) => {
	const modal = useAppSelector(getDialogState('EditStatus'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditStatus', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.data ? 'Редактировать статус' : 'Добавить статус'}
			headerActions={<DeleteComponent {...(modal?.context as Context)} submit={submit} />}
			body={<Form {...context} submit={submit} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const DeleteComponent: FC<{ data?: IStatusForm; submit: (data: IStatusForm) => void }> = ({ data, submit }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const deleteHandler = async () => {
		if (!data) return
		if (data.id) submit({ ...data, status: 'deleted' })
		else submit({ ...data, status: undefined })
		dispatch(changeDialogIsOpen({ variant: 'EditStatus', isOpen: false }))
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
			confirmText={`Вы уверены, что хотите удалить статус "${data?.label}"?`}
		/>
	)
}
