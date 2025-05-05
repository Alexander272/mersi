import { FC } from 'react'
import {
	Button,
	CircularProgress,
	Divider,
	FormControl,
	FormControlLabel,
	IconButton,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	Switch,
	TextField,
	useTheme,
} from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { ICreateFormField, ICreateFormFieldDTO } from '../../types/create'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Fallback } from '@/components/Fallback/Fallback'
import {
	useCreateFieldToCreateFormMutation,
	useDeleteFieldToCreateFormMutation,
	useUpdateFieldToCreateFormMutation,
} from '../../formApiSlice'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'

type Context = { item?: ICreateFormField }

export const FieldDialog = () => {
	const modal = useAppSelector(getDialogState('CreateFormField'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.item?.stepName ? 'Редактировать поле' : 'Добавить поле'}
			headerActions={context?.item?.id ? <DeleteComponent {...(modal?.context as Context)} /> : undefined}
			body={<Form {...context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const defaultValues: ICreateFormFieldDTO = {
	id: '',
	sectionId: '',
	step: 0,
	stepName: '',
	field: '',
	fieldName: '',
	path: '',
	type: '',
	isRequired: false,
	position: 0,
}

const Form: FC<Context> = ({ item }) => {
	defaultValues.sectionId = item?.sectionId || ''
	defaultValues.step = item?.step || 0
	defaultValues.stepName = item?.stepName || ''
	defaultValues.position = item?.position || 0

	const dispatch = useAppDispatch()

	const [create, { isLoading: creating }] = useCreateFieldToCreateFormMutation()
	const [update, { isLoading: updating }] = useUpdateFieldToCreateFormMutation()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<ICreateFormFieldDTO>({
		values: item?.id ? item : defaultValues,
	})

	const closeHandler = () => dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: false }))

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		const newData: ICreateFormFieldDTO = {
			id: item?.id || '',
			sectionId: form.sectionId,
			step: form.step,
			stepName: form.stepName.trim(),
			field: form.field.trim(),
			fieldName: form.fieldName.trim(),
			path: form.path.trim(),
			type: form.type,
			isRequired: form.isRequired,
			position: form.position,
		}

		try {
			if (item?.id) {
				await update(newData).unwrap()
				toast.success('Поле обновлено')
			} else {
				await create(newData).unwrap()
				toast.success('Поле создано')
			}
			dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name={'fieldName'}
				render={({ field }) => <TextField {...field} label={'Название поля'} fullWidth />}
			/>
			<Controller
				control={control}
				name={'field'}
				render={({ field }) => <TextField {...field} label={'Название поля (в объекте)'} fullWidth />}
			/>
			<Controller
				control={control}
				name={'path'}
				//TODO Может как-нибудь по другому это все обозвать
				render={({ field }) => <TextField {...field} label={'Группа (таблица)'} fullWidth />}
			/>
			<FormControl>
				<InputLabel id={'type'}>Тип поля</InputLabel>
				<Controller
					control={control}
					name='type'
					render={({ field, fieldState: { error } }) => (
						<Select labelId={'type'} label={'Тип поля'} error={Boolean(error)} {...field}>
							<MenuItem value='' disabled>
								Выберите тип поля
							</MenuItem>
							<MenuItem value='text'>Текст</MenuItem>
							<MenuItem value='number'>Число</MenuItem>
							<MenuItem value='date'>Дата</MenuItem>
							<MenuItem value='file'>Файл</MenuItem>
							<MenuItem value='list'>Список</MenuItem>
							{/* //TODO наверное стоит что-нибудь более понятное написать */}
							<MenuItem value='autocomplete'>Текст с авто дополнениями</MenuItem>
						</Select>
					)}
				/>
			</FormControl>
			<Controller
				control={control}
				name={'isRequired'}
				render={({ field }) => (
					<FormControlLabel
						control={<Switch checked={field.value || false} {...field} />}
						label={<>Поле {!field.value && 'не '}обязательно</>}
						sx={{ userSelect: 'none' }}
					/>
				)}
			/>

			<Divider sx={{ width: '50%', alignSelf: 'center' }} />
			<Stack spacing={2} direction={'row'}>
				<Button type='submit' variant='contained' fullWidth>
					Сохранить
				</Button>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отмена
				</Button>
			</Stack>

			{creating || updating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} mt={'0!important'} />
			) : null}
		</Stack>
	)
}

const DeleteComponent: FC<Context> = ({ item }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const [remove, { isLoading }] = useDeleteFieldToCreateFormMutation()

	const deleteHandler = async () => {
		if (!item?.id) return
		try {
			await remove(item.id).unwrap()
			toast.success('Поле удалено')
			dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Confirm
			buttonComponent={
				<IconButton size='large' sx={{ fill: '#505050', ':hover': { fill: palette.error.main } }}>
					{isLoading ? (
						<CircularProgress color='error' size={16} />
					) : (
						<DeleteIcon fontSize={16} fill={'inherit'} transition={'all 0.2s ease-in-out'} />
					)}
				</IconButton>
			}
			width='56px'
			onClick={deleteHandler}
			confirmText={`Вы уверены, что хотите удалить поле "${item?.fieldName}"?`}
		/>
	)
}
