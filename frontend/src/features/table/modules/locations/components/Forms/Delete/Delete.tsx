import { FC } from 'react'
import { Button, Divider, Stack, Typography } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { useAppDispatch } from '@/hooks/redux'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useDeleteLocationMutation, useGetLastLocationQuery } from '../../../locationsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Confirm } from '@/components/Confirm/Confirm'

type Props = {
	id: string
}

export const Delete: FC<Props> = ({ id }) => {
	const dispatch = useAppDispatch()

	const { data: instrument, isFetching } = useGetInstrumentByIdQuery(id || '', { skip: !id })
	const { data, isFetching: isFetchLast } = useGetLastLocationQuery(id || '', { skip: !id })
	const [remove, { isLoading }] = useDeleteLocationMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'DeleteLocation', isOpen: false }))
	}

	const deleteHandler = async () => {
		if (!data?.data.id) return

		try {
			await remove(data?.data.id).unwrap()
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack position={'relative'} mt={-2.5}>
			{isFetching || isFetchLast || isLoading ? <BoxFallback /> : null}

			<Typography fontSize={'1.2rem'} fontWeight={'bold'} textAlign={'center'}>
				{instrument?.data.name} ({instrument?.data.factoryNumber})
			</Typography>

			<Typography mt={2} mb={3} fontSize={'1.1rem'}>
				Удалить перемещение?
			</Typography>

			<Divider sx={{ width: '50%', alignSelf: 'center' }} />
			<Stack spacing={2} direction={'row'} mt={2}>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отменить
				</Button>
				<Confirm
					width='100%'
					onClick={deleteHandler}
					confirmText='Вы уверены, что хотите удалить перемещение?'
					buttonComponent={
						<Button variant='contained' fullWidth>
							Да
						</Button>
					}
				/>
			</Stack>
		</Stack>
	)
}
