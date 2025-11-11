import { FC } from 'react'
import { Button, Divider, FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'

import type { IHistoryForm } from '../../types/historyTypes'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'

const defaultValues: IHistoryForm = {
	id: '',
	sectionId: '',
	group: '',
	label: '',
	position: 1,
	status: 'new',
}

type Props = {
	data?: IHistoryForm
	submit: (data: IHistoryForm) => void
}

export const HistoryForm: FC<Props> = ({ data, submit }) => {
	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IHistoryForm>({
		defaultValues: data ? { ...data, status: data.status == 'new' ? 'new' : 'updated' } : defaultValues,
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditHistory', isOpen: false }))
	}

	const saveHandler = handleSubmit(form => {
		console.log('save', form, dirtyFields)
		if (Object.keys(dirtyFields).length) submit(form)
		closeHandler()
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name='label'
				render={({ field }) => <TextField {...field} label={'Название'} fullWidth />}
			/>

			<FormControl>
				<InputLabel id={'groups'}>Категория</InputLabel>
				<Controller
					control={control}
					name='group'
					render={({ field, fieldState: { error } }) => (
						<Select labelId={'groups'} label={'Категория'} error={Boolean(error)} {...field}>
							<MenuItem value='' disabled>
								Выберите категорию
							</MenuItem>
							<MenuItem value='verifications'>Поверки</MenuItem>
							<MenuItem value='locations'>Перемещения</MenuItem>
							<MenuItem value='repair'>Ремонт</MenuItem>
							<MenuItem value='preservation'>Консервация</MenuItem>
							<MenuItem value='save'>Хранение</MenuItem>
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
