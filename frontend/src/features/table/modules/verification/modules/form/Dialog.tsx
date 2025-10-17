import { FC } from 'react'
import { IconButton, useTheme } from '@mui/material'

import type { IVerificationFieldForm } from '../types/verificationFields'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { VerificationFieldsForm } from './Form'

type Context = { data?: IVerificationFieldForm }

type Props = {
	submit: (data: IVerificationFieldForm) => void
}

export const VerificationFieldsDialog: FC<Props> = ({ submit }) => {
	const modal = useAppSelector(getDialogState('EditVerificationFields'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditVerificationFields', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.data ? 'Редактировать поле формы' : 'Добавить поле формы'}
			headerActions={<DeleteComponent {...(modal?.context as Context)} submit={submit} />}
			body={<VerificationFieldsForm {...context} submit={submit} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const DeleteComponent: FC<{ data?: IVerificationFieldForm; submit: (data: IVerificationFieldForm) => void }> = ({
	data,
	submit,
}) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const deleteHandler = async () => {
		if (!data) return
		if (data.id) submit({ ...data, status: 'deleted' })
		else submit({ ...data, status: undefined })
		dispatch(changeDialogIsOpen({ variant: 'EditVerificationFields', isOpen: false }))
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
			confirmText={`Вы уверены, что хотите удалить поле формы "${data?.label}"?`}
		/>
	)
}
