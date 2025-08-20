import { FC } from 'react'
import { Button, Divider, Stack, Typography } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { useAppDispatch } from '@/hooks/redux'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useForcedReceivingMutation } from '../../../locationsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

type Props = {
	id: string
}

export const Forced: FC<Props> = ({ id }) => {
	const dispatch = useAppDispatch()

	const { data: instrument, isFetching } = useGetInstrumentByIdQuery(id || '', { skip: !id })
	const [receipt, { isLoading }] = useForcedReceivingMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Forced', isOpen: false }))
	}

	const receiveHandler = async () => {
		if (!id) return
		try {
			await receipt({ instrumentId: id }).unwrap()
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack position={'relative'} mt={-2.5}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Typography width={'100%'}>
				Подтвердите проставление отметки о получении инструмента «{instrument?.data.name}» (
				{instrument?.data.factoryNumber})
			</Typography>

			<Divider sx={{ width: '50%', alignSelf: 'center', mt: 3 }} />
			<Stack spacing={2} direction={'row'} mt={2}>
				<Button onClick={receiveHandler} variant='outlined' fullWidth>
					Подтвердить
				</Button>
			</Stack>
		</Stack>
	)
}
