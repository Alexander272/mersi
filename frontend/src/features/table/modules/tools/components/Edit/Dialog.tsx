import { FC } from 'react'
import { CircularProgress, IconButton, useTheme } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IToolsMenu } from '../../types/toolsMenu'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useDeleteToolsMenuMutation } from '../../toolsMenuApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { EditToolsForm } from './Form'

type Context = { data?: IToolsMenu; section: string }

export const EditToolsDialog = () => {
	const modal = useAppSelector(getDialogState('EditToolsMenu'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditToolsMenu', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.data ? 'Редактировать меню' : 'Добавить меню'}
			headerActions={<DeleteComponent {...(modal?.context as Context)} />}
			body={<EditToolsForm {...context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const DeleteComponent: FC<Context> = ({ data }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const [remove, { isLoading }] = useDeleteToolsMenuMutation()

	const deleteHandler = async () => {
		if (!data?.id) return

		try {
			await remove(data.id).unwrap()
			toast.success('Элемент меню "Инструменты" удален')
			dispatch(changeDialogIsOpen({ variant: 'EditToolsMenu', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Confirm
			buttonComponent={
				<IconButton size='large' sx={{ fill: '#505050', ':hover': { fill: palette.error.main } }}>
					{isLoading ? (
						<CircularProgress color='error' size={16} />
					) : (
						<DeleteIcon fontSize={16} fill={'inherit'} transition={'all 0.2s ease-in-out'} />
					)}
				</IconButton>
			}
			width='56px'
			onClick={deleteHandler}
			confirmText={`Вы уверены, что хотите удалить пункт меню "${data?.label}"?`}
		/>
	)
}
