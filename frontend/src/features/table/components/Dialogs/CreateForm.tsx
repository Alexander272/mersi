import { FC, useEffect, useState } from 'react'
import { Button, Divider, Stack, Tooltip } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ISiForm } from '../../types/si'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetCreateFormStepsQuery } from '@/features/sections/modules/form/formApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Step, Stepper } from '@/components/Stepper/Stepper'
import { localKeys } from '../../constants/storage'
import { useCreateSiMutation } from '../../siApiSlice'
import { Form as FormFields } from '../Form/Form'
import { useGetSI } from '../../hooks/getSI'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'

type Props = {
	id: string
}

export const CreateForm: FC<Props> = () => {
	const [activeStep, setActiveStep] = useState(0)
	const [steps, setSteps] = useState<Step[]>([])

	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(section?.id || '', { skip: !section?.id })
	const { data: si } = useGetSI()
	const [create, { isLoading }] = useCreateSiMutation()

	useEffect(() => {
		if (!data) return
		const newSteps = data.data.map(d => ({ id: d.step.toString(), label: d.stepName }))
		setSteps(newSteps)
	}, [data])

	const methods = useForm<ISiForm>({ values: JSON.parse(localStorage.getItem(localKeys.form) || '{}') })

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

		form.instrument.sectionId = section?.id || ''
		form.instrument.position = (si?.total || 0) + 1
		if (form.verification.verificationDate != 0 && form.instrument.interVerificationInterval != '') {
			form.verification.nextVerificationDate = dayjs(form.verification.verificationDate * 1000)
				.add(+form.instrument.interVerificationInterval, 'month')
				.unix()
		}

		try {
			await create(form).unwrap()
			toast.success('Данные добавлены')
			localStorage.removeItem(localKeys.form)
			methods.reset()
			setActiveStep(0)
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = () => {
		console.log('delete')
		localStorage.removeItem(localKeys.form)
		methods.reset()
	}

	return (
		<Stack position={'relative'} mt={-2}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack direction={'row'} width={'100%'} alignItems={'center'} mb={1.5}>
				{steps.length > 1 ? <Stepper steps={steps} active={activeStep} sx={{ width: '100%' }} /> : null}

				<Tooltip title='Очистить' enterDelay={600}>
					<span>
						<Button variant='outlined' color='inherit' onClick={deleteHandler}>
							<RefreshIcon fontSize={18} />
						</Button>
					</span>
				</Tooltip>
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
