import { FC, useRef } from 'react'
import { Box, Button, CircularProgress, MenuItem, Select, Stack, Tooltip, Typography, useTheme } from '@mui/material'
import { Controller, useFieldArray, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IDepartment } from '../types/departments'
import type { IResponsible } from '../types/responsible'
import { useGetUserByAccessQuery, useSyncUsersMutation } from '@/features/user/usersApiSlice'
import { useChangeResponsibleMutation, useGetResponsibleQuery } from '../responsibleApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'
import { SyncIcon } from '@/components/Icons/SyncIcon'
import { QuestionIcon } from '@/components/Icons/QuestionIcon'

type Props = {
	department?: IDepartment
}

export const Responsible: FC<Props> = ({ department }) => {
	const { palette } = useTheme()
	const defaultValue = { id: '', departmentId: department?.id || '', userId: '' }

	const { data, isFetching } = useGetResponsibleQuery({ department: department?.id || '' }, { skip: !department })
	const { data: users, isFetching: usersIsFetching } = useGetUserByAccessQuery(null)

	const [change] = useChangeResponsibleMutation()

	const [sync, { isLoading }] = useSyncUsersMutation()

	const deleted = useRef<string[]>([])

	const {
		control,
		setValue,
		handleSubmit,
		getValues,
		formState: { dirtyFields },
	} = useForm<{ values: IResponsible[] }>({
		values: { values: data?.data?.length ? data.data : [defaultValue] },
	})
	const { fields, append, remove } = useFieldArray({ control, name: 'values' })

	const addHandler = () => {
		append(defaultValue)
	}
	const removeHandler = (index: number) => {
		const values = getValues()
		const del = values.values[index].id
		deleted.current.push(del)

		if (index == 0 && fields.length == 1) setValue('values.0', defaultValue)
		else remove(index)
	}

	const syncHandler = async () => {
		try {
			await sync(null).unwrap()
			toast.success('Пользователи синхронизированы')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)

		const created: IResponsible[] = []
		const updated: IResponsible[] = []

		form.values.forEach(item => {
			if (item.userId) {
				if (item.id) {
					updated.push(item)
				} else {
					created.push(item)
				}
			}
		})

		try {
			await change({ new: created, updated, deleted: deleted.current }).unwrap()
			toast.success('Ответственные сохранены')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}

		console.log('deleted', deleted.current)
	})

	//TODO: лучше сделать автокмоплит с множественным выбором и автосохранением
	return (
		<Stack direction={'row'} alignItems={'flex-start'} position={'relative'} width={'100%'} mb={1}>
			{isFetching || usersIsFetching ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			<Stack spacing={2} mt={1} mr={1}>
				<Stack direction={'row'}>
					<Typography>Ответственный сотрудник(и): </Typography>
					<Tooltip title='Выбор пользователя, который может получать инструменты'>
						<Box
							sx={{
								cursor: 'pointer',
								mt: -0.5,
								width: 20,
								height: 20,
								display: 'flex',
								borderRadius: '50%',
								background: '#eff3fc',
								justifyContent: 'center',
								alignItems: 'center',
							}}
						>
							<QuestionIcon fontSize={12} />
						</Box>
					</Tooltip>
				</Stack>

				<Button
					onClick={saveHandler}
					sx={{
						textTransform: 'inherit',
						minHeight: 40,
						backgroundColor: '#9ab2ef29',
						':hover': { backgroundColor: '#9ab2ef52' },
					}}
				>
					Сохранить
				</Button>
			</Stack>

			<Stack spacing={1} width={'54%'}>
				{/* {fields.map((field, index) => (
					<Controller
						key={field.id}
						control={control}
						name={`values.${index}.ssoId`}
						render={({ field }) => {
							const user = users?.data.find(user => user.ssoId === field.value)
							return (
								<Autocomplete
									{...field}
									fullWidth
									value={user ? `${user.lastName} ${user.firstName}` : ''}
									onChange={(_e, value) => {
										const user = users?.data.find(
											user => `${user.lastName} ${user.firstName}` === value
										)
										field.onChange(user?.ssoId || '')
									}}
									options={users?.data.map(user => `${user.lastName} ${user.firstName}`) || []}
									renderInput={params => <TextField {...params} />}
									noOptionsText={'Ничего не найдено'}
								/>
							)
						}}
					/>
				))} */}
				{fields.map((field, index) => (
					<Controller
						key={field.id}
						control={control}
						name={`values.${index}.userId`}
						render={({ field }) => (
							<Select {...field}>
								{users?.data.map(user => (
									<MenuItem key={user.id} value={user.ssoId}>
										{user.lastName} {user.firstName}
									</MenuItem>
								))}
							</Select>
						)}
					/>
				))}
			</Stack>

			<Stack ml={2} spacing={1}>
				{fields.map((field, index) => (
					<Button
						key={field.id}
						onClick={() => removeHandler(index)}
						variant={'outlined'}
						color='error'
						sx={{ minWidth: 56, height: 40 }}
					>
						<FileDeleteIcon fontSize={16} fill={palette.error.main} />
					</Button>
				))}
			</Stack>
			<Stack justifyContent={'space-between'} ml={2} spacing={1}>
				<Tooltip title={'Синхронизировать список пользователей с системой авторизации'}>
					<Button onClick={syncHandler} sx={{ minWidth: 56, height: 40 }}>
						{isLoading ? <CircularProgress size={18} /> : <SyncIcon fontSize={18} />}
					</Button>
				</Tooltip>

				<Button onClick={addHandler} variant={'outlined'} sx={{ minWidth: 56, height: 40 }}>
					<PlusIcon fontSize={16} fill={palette.primary.main} />
				</Button>
			</Stack>
		</Stack>
	)
}
