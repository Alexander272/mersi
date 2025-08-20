import { FC } from 'react'
import { Button, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IEmployee } from '../types/employee'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useCreateEmployeeMutation, useGetEmployeeByIdQuery, useUpdateEmployeeMutation } from '../employeesApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Fallback } from '@/components/Fallback/Fallback'

type Context = { id?: string; department?: string }

export const EmployeeDialog = () => {
	const modal = useAppSelector(getDialogState('Employee'))

	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Employee', isOpen: false }))
	}

	return (
		<Dialog
			title={(modal?.context as Context)?.id ? 'Редактировать сотрудника' : 'Добавить сотрудника'}
			body={<Form {...(modal?.context as Context)} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}

const defaultValues = {
	id: '',
	name: '',
	departmentId: '',
	notes: '',
}
export const Form: FC<Context> = ({ id, department }) => {
	const { data, isFetching } = useGetEmployeeByIdQuery(id || '', { skip: !id })

	const dispatch = useAppDispatch()

	const [create] = useCreateEmployeeMutation()
	const [update] = useUpdateEmployeeMutation()

	const { control, handleSubmit } = useForm<IEmployee>({
		values: data?.data || { ...defaultValues, departmentId: department || '' },
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Employee', isOpen: false }))
	}

	const saveHandler = handleSubmit(async form => {
		console.log('save', form)

		const newData = {
			...data?.data,
			id: data?.data.id || id || '',
			name: form.name.trim(),
			notes: form.notes.trim(),
			departmentId: form.departmentId,
		}

		try {
			if (id) await update(newData).unwrap()
			else await create(newData).unwrap()
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} paddingX={2} position={'relative'} spacing={2} onSubmit={saveHandler}>
			{isFetching ? <Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} /> : null}

			<Controller
				control={control}
				name='name'
				rules={{ required: true }}
				render={({ field, fieldState: { error } }) => (
					<TextField {...field} label='ФИО' error={Boolean(error)} />
				)}
			/>
			<Controller
				control={control}
				name='notes'
				render={({ field, fieldState: { error } }) => (
					<TextField {...field} label='Доп. информация' error={Boolean(error)} />
				)}
			/>

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
