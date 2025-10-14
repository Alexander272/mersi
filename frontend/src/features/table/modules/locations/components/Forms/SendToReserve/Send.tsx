import { FC } from 'react'
import { Button, Divider, Stack, Table, TableBody, TableCell, TableHead, TableRow, Typography } from '@mui/material'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ILocationDTO } from '../../../types/location'
import { NullDate } from '@/constants/defaultValues'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetSIQuery } from '@/features/table/siApiSlice'
import { useCreateSeveralLocationsMutation } from '../../../locationsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

type Props = {
	ids: string[]
}

export const SendToReserve: FC<Props> = ({ ids }) => {
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const filter = { field: 'id', value: ids?.join(','), fieldType: 'list' as const, compareType: 'in' as const }
	const { data, isFetching } = useGetSIQuery(
		{ section: section?.id || '', status: 'work', filters: [filter], size: 9999999, all: true },
		{ skip: !ids }
	)
	const [create, { isLoading }] = useCreateSeveralLocationsMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'SendToReserve', isOpen: false }))
	}

	const saveHandler = async () => {
		const date = dayjs().startOf('d').toISOString()
		const locations: ILocationDTO[] = []

		ids.forEach(id => {
			locations.push({
				instrumentId: id,
				department: '',
				person: '',
				dateOfIssue: date,
				dateOfReceiving: NullDate,
				needConfirm: true,
				status: 'moved',
			})
		})

		try {
			if (!locations.length) return
			const payload = await create(locations).unwrap()
			toast.success(payload.message)
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

			<Table>
				<TableHead>
					<TableRow>
						<TableCell>Наименование</TableCell>
						<TableCell>Заводской номер</TableCell>
						<TableCell>Место нахождения</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{data?.data.map(r => (
						<TableRow key={r.id}>
							<TableCell>{r.name}</TableCell>
							<TableCell>{r.factoryNumber}</TableCell>
							<TableCell>{r.place}</TableCell>
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
