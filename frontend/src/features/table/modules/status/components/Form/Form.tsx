import { FC } from 'react'
import { Button, Divider, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'

import type { IStatusForm } from '../../types/status'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'

const defaultValues: IStatusForm = {
	id: '',
	position: 1,
	value: '',
	label: '',
	status: 'new',
}

type Props = {
	data?: IStatusForm
	submit: (data: IStatusForm) => void
}

export const Form: FC<Props> = ({ data, submit }) => {
	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IStatusForm>({
		defaultValues: data ? { ...data, status: data.status == 'new' ? 'new' : 'updated' } : defaultValues,
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditStatus', isOpen: false }))
	}

	const saveHandler = handleSubmit(form => {
		console.log('save', form, dirtyFields)
		if (Object.keys(dirtyFields).length) submit(form)
		dispatch(changeDialogIsOpen({ variant: 'EditStatus', isOpen: false }))
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name='label'
				render={({ field }) => <TextField {...field} label={'Название поля'} fullWidth />}
			/>

			<Controller
				control={control}
				name='value'
				render={({ field }) => <TextField {...field} label={'Код статуса'} fullWidth />}
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
