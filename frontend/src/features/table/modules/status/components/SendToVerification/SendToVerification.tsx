import { FC } from 'react'
import {
	Button,
	Divider,
	IconButton,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
	Typography,
} from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { useGetSIQuery } from '@/features/table/siApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { useChangeStatusMutation } from '@/features/table/instrumentApiSlice'
import { toast } from 'react-toastify'
import { IFetchError } from '@/app/types/error'

type Context = string | string[]

export const SendToVerificationDialog: FC<{ title: string }> = ({ title }) => {
	const modal = useAppSelector(getDialogState('SendToVerification'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'SendToVerification', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={title}
			headerActions={
				<IconButton onClick={closeHandler} size='large' sx={{ fill: '#505050', mr: 2 }}>
					<TimesIcon fontSize={12} />
				</IconButton>
			}
			body={<Content ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='sm'
			fullWidth
		/>
	)
}

const Content: FC<{ ids: string[] }> = ({ ids }) => {
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const filter = { field: 'id', value: ids?.join(','), fieldType: 'list' as const, compareType: 'in' as const }
	const { data, isFetching } = useGetSIQuery(
		{ section: section?.id || '', status: 'work', filters: [filter], size: 9999999, all: true },
		{ skip: !ids },
	)
	const [send, { isLoading }] = useChangeStatusMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'SendToVerification', isOpen: false }))
	}

	const saveHandler = async () => {
		const dto: { id: string; status: string }[] = []
		data?.data.forEach(d => {
			dto.push({ id: d.id, status: 'checking' })
		})

		try {
			if (!dto.length) {
				toast.warn('Инструменты не выбраны')
				return
			}
			await send(dto).unwrap()
			toast.success('Инструменты отправлены на проверку')
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
			console.log(error)
		}
	}

	if (!ids?.length) return <Typography textAlign={'center'}>Инструменты не выбраны</Typography>
	return (
		<Stack position={'relative'} mt={-2.5}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Typography textAlign={'center'} fontWeight={'bold'} mb={2}>
				Подтвердите выбор инструментов
			</Typography>

			<Table size='small'>
				<TableHead>
					<TableRow>
						<TableCell>Наименование</TableCell>
						<TableCell>Тип</TableCell>
						<TableCell>Заводской номер</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{data?.data.map(r => (
						<TableRow key={r.id}>
							<TableCell>{r.name}</TableCell>
							<TableCell>{r.type}</TableCell>
							<TableCell>{r.factoryNumber}</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>

			<Divider sx={{ width: '50%', alignSelf: 'center', mt: 3 }} />
			<Stack spacing={2} direction={'row'} mt={2}>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отмена
				</Button>
				<Button onClick={saveHandler} variant='contained' fullWidth>
					Сохранить
				</Button>
			</Stack>
		</Stack>
	)
}
