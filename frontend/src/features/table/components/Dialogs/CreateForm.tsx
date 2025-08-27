import { FC, useEffect, useState } from 'react'
import { Box, Button, Divider, Stack, Tooltip } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ISiForm } from '../../types/si'
import { localKeys } from '../../constants/storage'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetSI } from '../../hooks/getSI'
import { useGetCreateFormStepsQuery } from '@/features/sections/modules/form/formApiSlice'
import { useCreateSiMutation } from '../../siApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Step, Stepper } from '@/components/Stepper/Stepper'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { Form as FormFields } from '../Form/Form'

type Props = {
	id: string
}

export const CreateForm: FC<Props> = () => {
	const [activeStep, setActiveStep] = useState(0)
	const [steps, setSteps] = useState<Step[]>([])

	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(
		{ section: section?.id || '', action: 'Create' },
		{ skip: !section?.id }
	)
	const { data: si } = useGetSI()
	const [create, { isLoading }] = useCreateSiMutation()

	useEffect(() => {
		if (!data) return
		const newSteps = data.data.map(d => ({ id: d.step.toString(), label: d.stepName }))
		setSteps(newSteps)
	}, [data])

	const methods = useForm<ISiForm>({
		defaultValues: {},
		values: JSON.parse(localStorage.getItem(localKeys.form) || '{}'),
	})

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
		form.instrument.name = form.instrument.name.trim()
		form.instrument.type = form.instrument.type?.trim()
		form.instrument.factoryNumber = form.instrument.factoryNumber?.trim()
		form.instrument.measurementLimits = form.instrument.measurementLimits?.trim()
		form.instrument.accuracy = form.instrument.accuracy?.trim()
		form.instrument.stateRegister = form.instrument.stateRegister?.trim()
		form.instrument.manufacturer = form.instrument.manufacturer?.trim()
		form.instrument.notes = form.instrument.notes?.trim()

		if (form.verification) {
			form.verification.notes = form.verification.notes?.trim()
			form.verification.registerLink = form.verification.registerLink?.trim()

			if (!form.instrument.interVerificationInterval) {
				form.verification.nextVerificationDate = dayjs(form.verification.verificationDate * 1000)
					.add(+form.instrument.interVerificationInterval, 'month')
					.unix()
			}
			if (form.verification.notVerified) {
				form.instrument.interVerificationInterval = 0
				form.verification.verificationDate = 0
				form.verification.nextVerificationDate = 0
			}
		}
		if (form?.verification.docs?.length) {
			form.verification.docs = form.verification.docs.filter(d => d.doc && d.doc != '')
		}

		if (form.location) {
			const isToReserve = form.location.isToReserve
			form.location.department = isToReserve ? '' : form.location.department
			form.location.person = isToReserve ? '' : form.location.person
			form.location.dateOfReceiving = !form.location.needConfirm || isToReserve ? form.location.dateOfIssue : 0
			form.location.status = form.location.needConfirm ? 'moved' : isToReserve ? 'reserve' : 'used'
		}

		try {
			await create(form).unwrap()
			toast.success('Данные добавлены')
			methods.reset({ instrument: {}, verification: undefined })
			localStorage.removeItem(localKeys.form)
			setActiveStep(0)
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = () => {
		console.log('delete')
		methods.reset({ instrument: {}, verification: undefined })
		localStorage.removeItem(localKeys.form)
		setActiveStep(0)
	}

	return (
		<Stack position={'relative'} mt={-2}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack direction={'row'} width={'100%'} alignItems={'center'} mb={1.5}>
				{steps.length > 1 ? <Stepper steps={steps} active={activeStep} sx={{ width: '100%' }} /> : null}

				<Tooltip title='Очистить' enterDelay={600}>
					<Box ml={'auto'}>
						<Button variant='outlined' color='inherit' onClick={deleteHandler}>
							<RefreshIcon fontSize={18} />
						</Button>
					</Box>
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
