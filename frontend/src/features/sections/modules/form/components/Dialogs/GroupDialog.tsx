import { FC } from 'react'
import { Button, Divider, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { ICreateFormStep } from '../../types/create'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Fallback } from '@/components/Fallback/Fallback'
import { useUpdateSeveralFieldsToCreateFormMutation } from '../../formApiSlice'

type Context = { item?: ICreateFormStep; section?: string }

export const GroupDialog = () => {
	const modal = useAppSelector(getDialogState('CreateFormGroup'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={context?.item?.stepName ? 'Редактировать группу' : 'Добавить группу'}
			body={<Form {...context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const Form: FC<Context> = ({ item, section }) => {
	const dispatch = useAppDispatch()
	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<{ name: string }>({ values: { name: item?.stepName ? item.stepName : '' } })

	const [update, { isLoading }] = useUpdateSeveralFieldsToCreateFormMutation()

	const closeHandler = () => dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: false }))

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		if (!item?.stepName) {
			const context = { item: { step: item?.step || 0, stepName: form.name, sectionId: section || '' } }
			dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: true, context }))
			dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: false }))
			return
		}

		const fields = item.fields.map(d => ({ ...d, step: item.step, stepName: form.name.trim() }))

		try {
			await update(fields).unwrap()
			dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name={'name'}
				render={({ field }) => <TextField {...field} label={'Название группы'} fullWidth />}
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

			{isLoading ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} mt={'0!important'} />
			) : null}
		</Stack>
	)
}
