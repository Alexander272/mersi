import { FC } from 'react'
import { Button, Divider, FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'

import type { IVerificationFieldForm } from '../types/verificationFields'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'

const defaultValues: IVerificationFieldForm = {
	id: '',
	sectionId: '',
	field: '',
	label: '',
	type: '',
	position: 1,
	width: 0,
	group: '',
	status: 'new',
}

type Props = {
	data?: IVerificationFieldForm
	submit: (data: IVerificationFieldForm) => void
}

export const VerificationFieldsForm: FC<Props> = ({ data, submit }) => {
	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IVerificationFieldForm>({
		defaultValues: data ? { ...data, status: data.status == 'new' ? 'new' : 'updated' } : defaultValues,
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditVerificationFields', isOpen: false }))
	}

	const saveHandler = handleSubmit(form => {
		console.log('save', form, dirtyFields)
		if (Object.keys(dirtyFields).length) submit(form)
		dispatch(changeDialogIsOpen({ variant: 'EditVerificationFields', isOpen: false }))
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name='label'
				render={({ field }) => <TextField {...field} label={'Название колонки'} fullWidth />}
			/>

			<Controller
				control={control}
				name='field'
				render={({ field }) => <TextField {...field} label={'Название поля (в объекте)'} fullWidth />}
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
							<MenuItem value='checkbox'>Флажок</MenuItem>
							<MenuItem value='file'>Файл</MenuItem>
							{/* //TODO наверное стоит что-нибудь более понятное написать */}
							<MenuItem value='files'>Группа файлов</MenuItem>
							<MenuItem value='list'>Список</MenuItem>
							{/* //TODO наверное стоит что-нибудь более понятное написать */}
							<MenuItem value='autocomplete'>Текст с авто дополнениями</MenuItem>
						</Select>
					)}
				/>
			</FormControl>

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
