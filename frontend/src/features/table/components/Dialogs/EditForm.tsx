import { FC, useEffect, useMemo, useState } from 'react'
import { Button, Divider, Stack, Tooltip, useTheme } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ISiForm } from '../../types/si'
import type { IVerificationDTO } from '../../modules/verification/types/verification'
import { NullDate } from '@/constants/defaultValues'
import { PermRules } from '@/constants/permissions'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetCreateFormStepsQuery } from '@/features/sections/modules/form/formApiSlice'
import { useDeleteSIMutation, useGetSIByIdQuery, useUpdateSIMutation } from '../../siApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Step, Stepper } from '@/components/Stepper/Stepper'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { Form as FormFields } from '../Form/Form'
import { AspectRatioIcon } from '@/components/Icons/AspectRatioIcon'

type Props = {
	id: string
}

export const EditForm: FC<Props> = ({ id }) => {
	const [activeStep, setActiveStep] = useState(0)
	const [steps, setSteps] = useState<Step[]>([])

	const { palette } = useTheme()

	const canEditDetails = useCheckPermission([
		PermRules.Repair.Write,
		PermRules.Preservation.Write,
		PermRules.WriteOff.Write,
		PermRules.TransferToDep.Write,
		PermRules.TransferToSave.Write,
	])

	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(
		{ section: section?.id || '', action: 'Update' },
		{ skip: !section?.id },
	)
	const { data: si } = useGetSIByIdQuery(id, { skip: !id })
	const [update, { isLoading }] = useUpdateSIMutation()
	const [remove] = useDeleteSIMutation()

	useEffect(() => {
		if (!data) return
		const newSteps = data.data.map(d => ({ id: d.step.toString(), label: d.stepName }))
		setSteps(newSteps)
	}, [data])

	const formValues = useMemo(() => {
		if (!si?.data) return undefined
		const data = { ...si.data }
		if (!data.verification && !data.instrument.interVerificationInterval) {
			data.verification = {
				notVerified: true,
				instrumentId: data.instrument.id,
				id: '',
				verificationDate: '',
				nextVerificationDate: '',
				registerLink: '',
				status: '',
				notes: '',
			} as IVerificationDTO
		}
		return data as ISiForm
	}, [si?.data])

	const methods = useForm<ISiForm>({ values: formValues })

	const openDetailsDialog = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditTableItem', isOpen: false }))
		dispatch(changeDialogIsOpen({ variant: 'UpdateTableDetails', isOpen: true, context: { id } }))
	}

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

		if (!form.instrument.dateOfReceipt) form.instrument.dateOfReceipt = NullDate
		form.instrument.name = form.instrument.name.trim()
		form.instrument.type = form.instrument.type?.trim()
		form.instrument.factoryNumber = form.instrument.factoryNumber?.trim()
		form.instrument.measurementLimits = form.instrument.measurementLimits?.trim()
		form.instrument.accuracy = form.instrument.accuracy?.trim()
		form.instrument.stateRegister = form.instrument.stateRegister?.trim()
		form.instrument.manufacturer = form.instrument.manufacturer?.trim()
		form.instrument.notes = form.instrument.notes?.trim()

		if (!form.verification && !form.instrument.interVerificationInterval) {
			form.verification = {
				notVerified: true,
				instrumentId: form.instrument.id,
				id: '',
				verificationDate: '',
				nextVerificationDate: '',
				registerLink: '',
				status: '',
				notes: '',
			} as IVerificationDTO
		}

		if (form.verification) {
			form.verification.notes = form.verification.notes?.trim()
			form.verification.registerLink = form.verification.registerLink?.trim()

			if (
				form.verification.verificationDate &&
				form.verification.verificationDate != '' &&
				form.instrument.interVerificationInterval
			) {
				form.verification.nextVerificationDate = dayjs(form.verification.verificationDate)
					.add(+form.instrument.interVerificationInterval, 'month')
					.toISOString()
			}

			if (form.verification.notVerified) {
				form.instrument.interVerificationInterval = 0
				form.verification.verificationDate = NullDate
				form.verification.nextVerificationDate = NullDate
			}
		}
		if (form?.verification?.docs?.length) {
			const docs = form.verification.docs.filter(d => d.doc && d.doc != '')
			form.verification.docs = docs
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

				{canEditDetails && (
					<Tooltip title='Редактировать сведения'>
						<Button variant='outlined' color='primary' onClick={openDetailsDialog} sx={{ mr: 1 }}>
							<AspectRatioIcon fontSize={20} fill={palette.primary.main} />
						</Button>
					</Tooltip>
				)}

				<Confirm
					width='64px'
					onClick={deleteHandler}
					buttonComponent={
						<Button variant='outlined' color='error'>
							<FileDeleteIcon fontSize={20} fill={palette.error.main} />
						</Button>
					}
					confirmText='Вы уверены, что хотите удалить инструмент?'
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
