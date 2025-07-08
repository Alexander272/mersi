import { FC, useEffect, useState } from 'react'
import { Button, Divider, Stack, useTheme } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ISiForm } from '../../types/si'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetCreateFormStepsQuery } from '@/features/sections/modules/form/formApiSlice'
import { useDeleteSIMutation, useGetSIByIdQuery, useUpdateSIMutation } from '../../siApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Step, Stepper } from '@/components/Stepper/Stepper'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { Form as FormFields } from '../Form/Form'

type Props = {
	id: string
}

export const EditForm: FC<Props> = ({ id }) => {
	const [activeStep, setActiveStep] = useState(0)
	const [steps, setSteps] = useState<Step[]>([])

	const { palette } = useTheme()

	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(section?.id || '', { skip: !section?.id })
	const { data: si } = useGetSIByIdQuery(id, { skip: !id })
	const [update, { isLoading }] = useUpdateSIMutation()
	const [remove] = useDeleteSIMutation()

	useEffect(() => {
		if (!data) return
		const newSteps = data.data.map(d => ({ id: d.step.toString(), label: d.stepName }))
		setSteps(newSteps)
	}, [data])

	const methods = useForm<ISiForm>({ values: si?.data as ISiForm })

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditTableItem', isOpen: false }))
	}
	const prevHandler = () => setActiveStep(prev => prev - 1)

	const saveHandler = methods.handleSubmit(async form => {
		console.log('save', form, methods.formState.dirtyFields)

		if (activeStep + 1 != data?.data.length) {
			setActiveStep(prev => (prev + 1) % (data?.data.length || 0))
			return
		}
		if (!Object.keys(methods.formState.dirtyFields).length) {
			closeHandler()
			return
		}

		//TODO при обновлении если нет файла (свидетельства о поверке), то на сервер передается массив с пустым объектом и это вызывает ошибки

		if (form.verification.verificationDate != 0 && form.instrument.interVerificationInterval != '') {
			form.verification.nextVerificationDate = dayjs(form.verification.verificationDate * 1000)
				.add(+form.instrument.interVerificationInterval, 'month')
				.unix()
		}

		try {
			await update(form).unwrap()
			toast.success('Данные обновлены')
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = async () => {
		if (!data) return

		try {
			await remove(id).unwrap()
			toast.success('Данные удалены')
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack position={'relative'} mt={-2}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack direction={'row'} width={'100%'} alignItems={'center'} mb={1.5}>
				{steps.length > 1 ? <Stepper steps={steps} active={activeStep} sx={{ width: '100%' }} /> : null}

				<Confirm
					width='64px'
					onClick={deleteHandler}
					buttonComponent={
						<Button variant='outlined' color='error'>
							<FileDeleteIcon fontSize={20} fill={palette.error.main} />
						</Button>
					}
					confirmText='Вы уверены, что хотите удалить реактив?'
				/>
			</Stack>

			<Stack mt={2} component={'form'} onSubmit={saveHandler}>
				<FormProvider {...methods}>
					<FormFields data={data?.data[activeStep].fields || []} instrumentId={si?.data.instrument.id} />
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
