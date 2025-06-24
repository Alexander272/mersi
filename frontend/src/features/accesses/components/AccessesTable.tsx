import { FC } from 'react'
import {
	Button,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	useTheme,
} from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IAccesses } from '../types/accesses'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Confirm } from '@/components/Confirm/Confirm'
import { Fallback } from '@/components/Fallback/Fallback'
import { CloseIcon } from '@/components/Icons/CloseIcon'
import { EditIcon } from '@/components/Icons/EditIcon'
import { AccessDialog } from './AccessDialog'
import { useDeleteAccessMutation, useGetAccessesQuery } from '../accessesApiSlice'

type Props = {
	realm: string
}

export const AccessesTable: FC<Props> = ({ realm }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetAccessesQuery(realm, { skip: !realm || realm == 'new' })
	const [remove, { isLoading }] = useDeleteAccessMutation()

	const deleteHandler = async (id: string) => {
		const access = data?.data.find(e => e.id === id)
		if (!access) return

		try {
			await remove(access.id).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	const updateHandler = (item: IAccesses) => {
		dispatch(changeDialogIsOpen({ variant: 'Access', isOpen: true, content: { item, realm } }))
	}

	const addHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Access', isOpen: true, content: { realm } }))
	}

	return (
		<TableContainer sx={{ position: 'relative', minHeight: 140, mt: 3, mb: 4 }}>
			{isFetching || isLoading ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} alignItems={'flex-start'} pt={3} />
			) : null}

			<AccessDialog />

			<Table size='small'>
				<TableHead>
					<TableRow>
						<TableCell width={'30%'}>ФИО</TableCell>
						<TableCell width={'25%'}>Роль</TableCell>
						<TableCell width={'30%'}>Дата добавления</TableCell>
						<TableCell width={'15%'}>Действия</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id}>
							<TableCell>
								{item.user.lastName} {item.user.firstName}
							</TableCell>
							<TableCell>{item.role.name}</TableCell>
							<TableCell>{new Date(item.created).toLocaleString('ru-RU')}</TableCell>
							<TableCell>
								<Stack direction={'row'}>
									<Button onClick={() => updateHandler(item)} sx={{ minWidth: 46 }}>
										<EditIcon fontSize={16} />
									</Button>

									<Confirm
										onClick={() => deleteHandler(item.id)}
										width='46px'
										confirmText='Вы уверены, что хотите удалить сотрудника?'
										buttonComponent={
											<Button sx={{ minWidth: 46 }}>
												<CloseIcon fontSize={20} fill={palette.error.main} />
											</Button>
										}
									/>
								</Stack>
							</TableCell>
						</TableRow>
					))}

					<TableRow>
						<TableCell colSpan={4} align='center'>
							<Button
								onClick={addHandler}
								sx={{
									textTransform: 'inherit',
									width: 300,
									backgroundColor: '#9ab2ef29',
									':hover': { backgroundColor: '#9ab2ef52' },
								}}
							>
								Добавить
							</Button>
						</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</TableContainer>
	)
}
