import { FC, useEffect } from 'react'
import {
	Button,
	FormControl,
	FormControlLabel,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	Switch,
	TextField,
	Typography,
	useTheme,
} from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IRealm } from '../types/realm'
import { Fallback } from '@/components/Fallback/Fallback'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'
import {
	useCreateRealmMutation,
	useDeleteRealmMutation,
	useGetRealmByIdQuery,
	useUpdateRealmMutation,
} from '../realmsApiSlice'

type Props = {
	realm: string
	setRealm: (realm: string) => void
}

type Form = Omit<IRealm, 'id' | 'created'>

const defaultValues: Form = {
	name: '',
	realm: '',
	isActive: true,
	reserveChannel: '',
	expirationNotice: false,
	locationType: 'department',
	hasResponsible: true,
	needResponsible: true,
	needConfirmed: true,
}

export const RealmForm: FC<Props> = ({ realm, setRealm }) => {
	const { palette } = useTheme()

	const { data, isFetching } = useGetRealmByIdQuery(realm, { skip: !realm || realm == 'new' })

	const [create, { isLoading: isCreating }] = useCreateRealmMutation()
	const [update, { isLoading: isUpdating }] = useUpdateRealmMutation()
	const [remove, { isLoading: isDeleting }] = useDeleteRealmMutation()

	const {
		control,
		reset,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<Form>({ values: defaultValues })

	useEffect(() => {
		if (data && realm != 'new') reset(data.data)
		else reset(defaultValues)
	}, [data, reset, realm])

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		const newData = { ...form, id: data?.data.id || '' }
		try {
			if (realm == 'new') {
				const payload = await create(newData).unwrap()
				setRealm(payload.id)
				toast.success('Область создана')
			} else {
				await update(newData).unwrap()
				toast.success('Область обновлена')
			}
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = async () => {
		if (realm == 'new') return
		try {
			await remove(realm).unwrap()
			setRealm('new')
			toast.success('Область удалена')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack component={'form'} onSubmit={saveHandler} position={'relative'}>
			{isFetching || isDeleting || isUpdating || isCreating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'name'}
					render={({ field }) => <TextField {...field} label={'Область'} fullWidth />}
				/>

				<Controller
					control={control}
					name={'realm'}
					render={({ field }) => <TextField {...field} label={'Короткое название (код, url)'} fullWidth />}
				/>
			</Stack>

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'reserveChannel'}
					render={({ field }) => <TextField {...field} label={'Канал'} fullWidth />}
				/>

				<FormControl fullWidth>
					<InputLabel id='locationType'>Тип местоположения</InputLabel>
					<Controller
						control={control}
						name='locationType'
						render={({ field }) => (
							<Select labelId='locationType' {...field} label={'Тип местоположения'}>
								<MenuItem value='department'>Подразделение</MenuItem>
								<MenuItem value='place'>Место нахождения</MenuItem>
							</Select>
						)}
					/>
				</FormControl>
			</Stack>

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'needConfirmed'}
					render={({ field }) => (
						<FormControlLabel
							control={<Switch checked={field.value} {...field} />}
							label={
								<>
									Подтверждения инструментов
									<br />
									{field.value ? 'используются' : 'не используются'}
								</>
							}
						/>
					)}
				/>

				<Controller
					control={control}
					name='hasResponsible'
					render={({ field }) => (
						<FormControlLabel
							control={<Switch checked={field.value} {...field} />}
							label={field.value ? 'Есть сотрудники' : 'Нет сотрудников'}
						/>
					)}
				/>

				<Controller
					control={control}
					name='needResponsible'
					render={({ field }) => (
						<FormControlLabel
							control={<Switch checked={field.value} {...field} />}
							label={`${field.value ? 'Нужен' : 'Не нужен'} ответственный сотрудник`}
						/>
					)}
				/>
			</Stack>

			<Stack direction={'row'} flexGrow={1} alignItems={'center'} justifyContent={'space-between'}>
				<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
					<Controller
						control={control}
						name={'isActive'}
						render={({ field }) => (
							<FormControlLabel
								control={<Switch checked={field.value} {...field} />}
								label={field.value ? 'Активна' : 'Не активна'}
							/>
						)}
					/>

					<Controller
						control={control}
						name={'expirationNotice'}
						render={({ field }) => (
							<FormControlLabel
								control={<Switch checked={field.value} {...field} />}
								label={`Уведомления о поверках ${field.value ? 'включены' : 'выключены'}`}
							/>
						)}
					/>
				</Stack>
				<Stack direction={'row'} flexGrow={1} spacing={2} mb={1} justifyContent={'flex-end'} height={38}>
					<Button
						variant='outlined'
						type='submit'
						disabled={!Object.keys(dirtyFields).length}
						sx={{ minWidth: 56 }}
					>
						<SaveIcon
							fontSize={18}
							fill={!Object.keys(dirtyFields).length ? palette.action.disabled : palette.primary.main}
						/>
					</Button>

					<Confirm
						onClick={deleteHandler}
						buttonComponent={
							<Button
								variant='outlined'
								color='error'
								disabled={realm == 'new'}
								sx={{ minWidth: 56, height: '100%' }}
							>
								<FileDeleteIcon
									fontSize={20}
									fill={realm == 'new' ? palette.action.disabled : palette.error.main}
								/>
							</Button>
						}
						confirmText='Вы уверены, что хотите удалить область?'
						width='56'
					/>
				</Stack>
			</Stack>

			{realm != 'new' && data?.data ? (
				<Typography>Дата создания области: {new Date(data.data.created).toLocaleString('ru-RU')}</Typography>
			) : null}
		</Stack>
	)
}
