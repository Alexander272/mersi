import type { FC } from 'react'
import { Button, Stack, useTheme } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { ITransferToDepartment, ITransferToDepartmentDTO } from '../../types/department'
import { useAppDispatch } from '@/hooks/redux'
import { useDeleteTransferToDepartmentMutation, useUpdateTransferToDepartmentMutation } from '../../transferApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { Inputs } from './Inputs'

type Props = {
	data: ITransferToDepartment
}

export const EditTransferToDepItem: FC<Props> = ({ data }) => {
	const dispatch = useAppDispatch()
	const { palette } = useTheme()

	const methods = useForm<ITransferToDepartmentDTO>({
		values: { ...data, instrumentId: data.instrumentId || '' },
	})

	const [update, { isLoading }] = useUpdateTransferToDepartmentMutation()
	const [remove, { isLoading: isDeleting }] = useDeleteTransferToDepartmentMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'UpdateTableDetails', isOpen: false }))
	}

	const submitHandler = methods.handleSubmit(async form => {
		try {
			await update(form).unwrap()
			closeHandler()
		} catch (error) {
			toast.error((error as IFetchError).data.message, { autoClose: false })
		}
	})

	const deleteHandler = async () => {
		try {
			await remove({ id: data.id, instrumentId: data.instrumentId }).unwrap()
			closeHandler()
		} catch (error) {
			toast.error((error as IFetchError).data.message, { autoClose: false })
		}
	}

	return (
		<Stack position={'relative'} my={2} mx={3}>
			{isLoading || isDeleting ? <BoxFallback /> : null}

			<Stack component={'form'} onSubmit={submitHandler}>
				<FormProvider {...methods}>
					<Inputs instrumentId={data.instrumentId} />
				</FormProvider>

				<Stack direction={'row'} justifyContent={'center'} spacing={2}>
					<Button variant='outlined' type='submit' sx={{ textTransform: 'inherit', width: 300 }}>
						Сохранить
					</Button>

					<Confirm
						onClick={deleteHandler}
						confirmText='Вы действительно хотите удалить запись о передаче?'
						buttonComponent={
							<Button color='error' sx={{ height: 37 }}>
								<DeleteIcon fontSize={18} fill={palette.error.light} />
							</Button>
						}
					/>
				</Stack>
			</Stack>
		</Stack>
	)
}
