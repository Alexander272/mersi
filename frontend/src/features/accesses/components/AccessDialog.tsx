import { FC } from 'react'
import {
	Autocomplete,
	Button,
	CircularProgress,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	TextField,
	Typography,
} from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IAccesses, IAccessesDTO } from '../types/accesses'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { useGetRolesQuery } from '@/features/user/roleApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { useGetUsersByRealmQuery, useSyncUsersMutation } from '@/features/user/usersApiSlice'
import { IUserData } from '@/features/user/types/user'
import { useCreateAccessMutation, useUpdateAccessMutation } from '../accessesApiSlice'
import { SyncIcon } from '@/components/Icons/SyncIcon'

type Context = { item?: IAccesses; realm?: string }

export const AccessDialog = () => {
	const modal = useAppSelector(getDialogState('Access'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Access', isOpen: false }))
	}

	return (
		<Dialog
			title={(modal?.content as Context)?.item ? 'Редактировать доступ' : 'Добавить доступ'}
			body={<Form {...(modal?.content as Context)} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

type Form = {
	id?: string
	user: IUserData | null
	roleId: string
}

const defaultValues: Form = {
	user: null,
	roleId: '',
}

export const Form: FC<Context> = ({ item, realm }) => {
	// const
	const dispatch = useAppDispatch()

	const { data: users, isFetching: usersIsFetching } = useGetUsersByRealmQuery(
		{ realm: realm || '' },
		{ skip: !realm }
	)
	const { data: roles, isFetching: rolesIsFetching } = useGetRolesQuery(null)

	const [create, { isLoading: creating }] = useCreateAccessMutation()
	const [update, { isLoading: updating }] = useUpdateAccessMutation()

	const [sync, { isLoading: isSync }] = useSyncUsersMutation()

	const { control, handleSubmit } = useForm<Form>({
		values: item ? { id: item.id, user: item.user, roleId: item.role.id } : defaultValues,
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Access', isOpen: false }))
	}

	const saveHandler = handleSubmit(async form => {
		console.log('save', form)

		const dto: IAccessesDTO = {
			id: item?.id || '',
			realmId: realm || '',
			userId: form.user?.id || '',
			roleId: form.roleId,
		}

		try {
			if (item?.id) await update(dto).unwrap()
			else await create(dto).unwrap()
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const syncHandler = async () => {
		try {
			await sync(null).unwrap()
			toast.success('Пользователи синхронизированы')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			{rolesIsFetching || usersIsFetching || creating || updating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			<Button onClick={syncHandler} color='inherit' sx={{ textTransform: 'inherit' }}>
				Синхронизировать пользователей
				{isSync ? <CircularProgress size={18} sx={{ ml: 2 }} /> : <SyncIcon fontSize={18} ml={2} />}
			</Button>

			<Controller
				name='user'
				control={control}
				render={({ field: { onChange, ...props } }) => (
					<Autocomplete
						options={users?.data || []}
						getOptionLabel={option => `${option.lastName} ${option.firstName} (${option.email})`}
						renderOption={(props, option) => {
							return (
								<li {...props} key={option.id}>
									<Stack>
										<Typography>{`${option.lastName} ${option.firstName}`}</Typography>
										<Typography variant='body2' color='text.secondary'>
											{option.email}
										</Typography>
									</Stack>
								</li>
							)
						}}
						isOptionEqualToValue={(option, value) => option.id === value.id}
						renderInput={params => <TextField {...params} label='Пользователь' variant='outlined' />}
						onChange={(_e, data) => onChange(data)}
						{...props}
					/>
				)}
			/>

			<FormControl>
				<InputLabel id={'role'}>Роль</InputLabel>
				<Controller
					control={control}
					name='roleId'
					render={({ field, fieldState: { error } }) => (
						<Select labelId={'role'} label={'Роль'} error={Boolean(error)} {...field}>
							<MenuItem value='' disabled>
								Выберите роль
							</MenuItem>

							{roles?.data.map(o => (
								<MenuItem key={o.id} value={o.id}>
									{o.name} {o.description ? `(${o.description})` : ''}
								</MenuItem>
							))}
						</Select>
					)}
				/>
			</FormControl>

			<Stack spacing={2} direction={'row'}>
				<Button type='submit' variant='contained' fullWidth>
					Сохранить
				</Button>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отмена
				</Button>
			</Stack>
		</Stack>
	)
}
