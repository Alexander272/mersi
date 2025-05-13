import { useEffect, useState } from 'react'
import { Button, Divider, Stack } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetCreateFormStepsQuery } from '@/features/sections/modules/form/formApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Step, Stepper } from '@/components/Stepper/Stepper'
import { localKeys } from '../../constants/storage'
import { useCreateSiMutation } from '../../siApiSlice'
import { Form as FormFields } from '../Form/Form'

export const CreateDialog = () => {
	const modal = useAppSelector(getDialogState('CreateTableItem'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: false }))
	}

	return (
		<Dialog
			title={'Добавить'}
			body={<Form />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}

const Form = () => {
	const [activeStep, setActiveStep] = useState(0)
	const [steps, setSteps] = useState<Step[]>([])

	// const section = useAppSelector(getSection) //TODO получать нормальное значение section
	const section = '46ba9e17-65c7-474b-8c47-7975ab4319d5'
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(section, { skip: !section })
	const [create, { isLoading }] = useCreateSiMutation()

	useEffect(() => {
		if (!data) return
		const newSteps = data.data.map(d => ({ id: d.step.toString(), label: d.stepName }))
		setSteps(newSteps)
	}, [data])

	const methods = useForm({ values: JSON.parse(localStorage.getItem(localKeys.form) || '{}') })

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'CreateTableItem', isOpen: false }))
	}
	const prevHandler = () => setActiveStep(prev => prev - 1)

	const saveHandler = methods.handleSubmit(async form => {
		console.log('save', form, methods.formState.dirtyFields)
		if (Object.keys(methods.formState.dirtyFields).length) {
			console.log('change local storage')
			localStorage.setItem(localKeys.form, JSON.stringify(form))
		}
		if (activeStep + 1 != data?.data.length) {
			setActiveStep(prev => (prev + 1) % (data?.data.length || 0))
			return
		}

		form.instrument.sectionId = section
		//TODO надо еще что-то делать с полем nextVerificationData
		try {
			await create(form).unwrap()
			toast.success('Новая позиция добавлена')
			localStorage.removeItem(localKeys.form)
			methods.reset()
			setActiveStep(0)
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack position={'relative'} mt={-2}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack direction={'row'} width={'100%'} alignItems={'center'} mb={1.5}>
				{steps.length > 1 ? <Stepper steps={steps} active={activeStep} sx={{ width: '100%' }} /> : null}

				{/* <Tooltip title='Удалить черновик'>
					<span>
						<Button variant='outlined' color='error' onClick={deleteHandler} disabled={!instrument}>
							<FileDeleteIcon
								fill={!instrument ? palette.action.disabled : palette.error.main}
								fontSize={22}
							/>
						</Button>
					</span>
				</Tooltip> */}
			</Stack>

			<Stack mt={2} component={'form'} onSubmit={saveHandler}>
				<FormProvider {...methods}>
					<FormFields data={data?.data[activeStep].fields || []} />
				</FormProvider>

				<Divider sx={{ width: '50%', alignSelf: 'center' }} />
				<Stack spacing={2} direction={'row'} mt={2}>
					{activeStep == 0 ? (
						<Button onClick={closeHandler} variant='outlined' fullWidth>
							Отмена
						</Button>
					) : (
						<Button onClick={prevHandler} variant='outlined' fullWidth>
							Назад
						</Button>
					)}

					<Button type='submit' variant='contained' fullWidth>
						{activeStep == steps.length - 1 ? 'Сохранить' : 'Далее'}
					</Button>
				</Stack>
			</Stack>
		</Stack>
	)
}
