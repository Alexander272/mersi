import { FC } from 'react'
import { Button, Divider, Stack, Typography } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { PermRules } from '@/constants/permissions'
import { useAppDispatch } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useReceivingMutation } from '../../../locationsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

type Props = {
	id: string
}

export const ReceiptOne: FC<Props> = ({ id }) => {
	const hasResWrite = useCheckPermission(PermRules.Reserve.Write)
	const dispatch = useAppDispatch()

	const { data: instrument, isFetching } = useGetInstrumentByIdQuery(id || '', { skip: !id })
	const [receipt, { isLoading }] = useReceivingMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Receive', isOpen: false }))
	}

	const receiveHandler = async () => {
		if (!id) return
		try {
			await receipt({ instrumentId: [id], status: hasResWrite ? 'used' : 'reserve' }).unwrap()
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
				Подтвердите получение инструмента «{instrument?.data.name}» ({instrument?.data.factoryNumber})
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
