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
import type { IColumn, IColumnDTO } from '../../types/columns'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetGroupedSectionsQuery } from '@/features/sections/sectionsApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Fallback } from '@/components/Fallback/Fallback'
import { Confirm } from '@/components/Confirm/Confirm'
import { DeleteIcon } from '@/components/Icons/DeleteIcon'
import { useCreateColumnMutation, useDeleteColumnMutation, useUpdateColumnMutation } from '../../columnsApiSlice'

type Context = { item?: IColumn; section?: string }

export const ColumnDialog = () => {
	const modal = useAppSelector(getDialogState('Columns'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.item?.id ? 'Редактировать колонку' : 'Добавить колонку'}
			headerActions={
				context?.item?.id && (context?.item?.children?.length || 0) == 0 ? (
					<DeleteComponent {...(modal?.context as Context)} />
				) : undefined
			}
			body={<Form {...context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const defaultValues: IColumnDTO = {
	id: '',
	sectionId: '',
	name: '',
	field: '',
	position: 0,
	width: 200,
	type: '',
	allowSort: true,
	allowFilter: true,
}

const Form: FC<Context> = ({ item }) => {
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetGroupedSectionsQuery(null)

	const [create, { isLoading: creating }] = useCreateColumnMutation()
	const [update, { isLoading: updating }] = useUpdateColumnMutation()

	const def = {
		...defaultValues,
		sectionId: item?.sectionId || '',
		parentId: item?.parentId,
		position: item?.position || 0,
	}
	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IColumnDTO>({
		values: item?.id ? item : def,
	})

	const closeHandler = () => dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: false }))

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		const newData: IColumnDTO = {
			id: item?.id || '',
			sectionId: form.sectionId,
			name: form.name.trim(),
			field: form.field.trim(),
			position: form.position,
			width: +(form.width || 0),
			type: form.type,
			allowSort: form.allowSort,
			allowFilter: form.allowFilter,
			parentId: form.parentId,
		}

		try {
			if (!item?.id) {
				await create(newData).unwrap()
				toast.success('Колонка создана')
			} else {
				await update(newData).unwrap()
				toast.success('Колонка обновлена')
			}
			dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			{isFetching || creating || updating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			{/* <IconButton size='large' sx={{ width: 40, height: 40, position: 'absolute', top: -10, right: 10 }}>
				<DeleteIcon fontSize={16} />
			</IconButton> */}

			<Controller
				control={control}
				name={'name'}
				render={({ field }) => <TextField {...field} label={'Название колонки'} fullWidth />}
			/>
			<Controller
				control={control}
				name={'field'}
				render={({ field }) => <TextField {...field} label={'Название поля (в объекте)'} fullWidth />}
			/>
			<FormControl>
				<InputLabel id={'type'}>Тип колонки</InputLabel>
				<Controller
					control={control}
					name='type'
					render={({ field, fieldState: { error } }) => (
						<Select labelId={'type'} label={'Тип колонки'} error={Boolean(error)} {...field}>
							<MenuItem value='' disabled>
								Выберите тип колонки
							</MenuItem>
							<MenuItem value='text'>Текст</MenuItem>
							<MenuItem value='number'>Число</MenuItem>
							<MenuItem value='date'>Дата</MenuItem>
							<MenuItem value='file'>Файл</MenuItem>
							<MenuItem value='list'>Список</MenuItem>
							{/* //TODO наверное стоит что-нибудь более понятное написать */}
							<MenuItem value='autocomplete'>Текст с авто дополнениями</MenuItem>
							<MenuItem value='parent'>Группа с вложенными колонками</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

			<Controller
				control={control}
				name={'width'}
				render={({ field }) => (
					<TextField
						{...field}
						label={'Название поля (в объекте)'}
						fullWidth
						slotProps={{
							htmlInput: {
								step: 1,
								min: 1,
								type: 'number',
							},
						}}
					/>
				)}
			/>

			<FormControl>
				<InputLabel id={'sectionId'}>Секция</InputLabel>
				<Controller
					control={control}
					name='sectionId'
					render={({ field, fieldState: { error } }) => (
						<Select labelId={'sectionId'} label={'Секция'} error={Boolean(error)} {...field}>
							<MenuItem value='' disabled>
								Выберите секцию
							</MenuItem>

							{data?.data.map(d => {
								return d.sections.map(item => (
									<MenuItem key={item.id} value={item.id}>
										{d.title} / {item.name}
									</MenuItem>
								))
							})}
						</Select>
					)}
				/>
			</FormControl>

			<Controller
				control={control}
				name={'allowSort'}
				render={({ field }) => (
					<FormControlLabel
						control={<Switch checked={field.value || false} {...field} />}
						label={<>Сортировка {!field.value && 'не '}доступна</>}
						sx={{ userSelect: 'none' }}
					/>
				)}
			/>
			<Controller
				control={control}
				name={'allowFilter'}
				render={({ field }) => (
					<FormControlLabel
						control={<Switch checked={field.value || false} {...field} />}
						label={<>Фильтрация {!field.value && 'не '}доступна</>}
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
		</Stack>
	)
}

const DeleteComponent: FC<Context> = ({ item }) => {
	const { palette } = useTheme()
	const dispatch = useAppDispatch()

	const [remove, { isLoading }] = useDeleteColumnMutation()

	const deleteHandler = async () => {
		if (!item?.id) return
		try {
			await remove(item.id).unwrap()
			toast.success('Колонка удалена')
			dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: false }))
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
			confirmText={`Вы уверены, что хотите удалить колонку "${item?.name}"?`}
		/>
	)
}
