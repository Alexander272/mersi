import { FC, useState } from 'react'
import { Button, Divider, IconButton, Stack, Typography, useTheme } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { IVerificationDTO } from '../../types/verification'
import { useAppDispatch } from '@/hooks/redux'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useCreateVerificationMutation } from '../../verificationApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { Inputs } from './Inputs'

const def = {
	verificationDate: dayjs().unix(),
	registerLink: '',
	status: 'work',
	notes: '',
}

type Props = {
	ids: string[]
}

export const Create: FC<Props> = ({ ids }) => {
	const [active, setActive] = useState(0)
	const { palette } = useTheme()

	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetInstrumentByIdQuery(ids?.length ? ids[active] : '', {
		skip: !ids?.length || !ids[active],
	})
	const [create, { isLoading }] = useCreateVerificationMutation()

	const methods = useForm<IVerificationDTO>({ defaultValues: def })

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'NewVerification', isOpen: false }))
	}

	const activeHandler = (type: 'prev' | 'next') => () => {
		if (type == 'prev') setActive(prev => prev - 1)
		else setActive(prev => prev + 1)
	}

	const saveHandler = methods.handleSubmit(async form => {
		if (!data) return
		console.log('save', form, methods.formState.dirtyFields)

		form.instrumentId = data.data.id
		form.nextVerificationDate = dayjs(form.verificationDate * 1000)
			.add(+(data.data.interVerificationInterval || 0), 'month')
			.unix()

		try {
			await create(form).unwrap()
			toast.success('Данные о поверке добавлены')
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	if (!ids?.length) return <Typography textAlign={'center'}>Инструменты не выбраны</Typography>
	return (
		<Stack position={'relative'} mt={-2.5}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack spacing={2} direction={'row'} paddingX={3}>
				{ids.length > 1 && (
					<IconButton onClick={activeHandler('prev')} disabled={active == 0}>
						<LeftArrowIcon fontSize={16} fill={active == 0 ? palette.action.disabled : palette.grey[900]} />
					</IconButton>
				)}

				<Typography textAlign={'center'} sx={{ width: '100%' }}>
					<Typography component={'span'} mr={2.5} fontSize={'1.3rem'} color='primary'>
						{ids.length > 1 ? `${active + 1}/${ids.length}` : ''}
					</Typography>
					{data?.data.name} ({data?.data.factoryNumber})
				</Typography>

				{ids.length > 1 && (
					<IconButton onClick={activeHandler('next')} disabled={active == ids.length - 1}>
						<LeftArrowIcon
							fontSize={16}
							transform={'rotate(180deg)'}
							fill={active == ids.length - 1 ? palette.action.disabled : palette.grey[900]}
						/>
					</IconButton>
				)}
			</Stack>

			<Stack mt={2} component={'form'} onSubmit={saveHandler}>
				<FormProvider {...methods}>
					<Inputs instrumentId={ids?.length ? ids[active] : ''} />
				</FormProvider>

				<Divider sx={{ width: '50%', alignSelf: 'center' }} />
				<Stack spacing={2} direction={'row'} mt={2}>
					<Button onClick={closeHandler} variant='outlined' fullWidth>
						Отмена
					</Button>
					<Button type='submit' variant='contained' fullWidth>
						Сохранить
					</Button>
				</Stack>
			</Stack>
		</Stack>
	)
}
