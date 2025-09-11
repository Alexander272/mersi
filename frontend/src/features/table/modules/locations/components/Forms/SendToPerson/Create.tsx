import { FC, useEffect, useRef, useState } from 'react'
import { Button, Divider, IconButton, Stack, Typography, useTheme } from '@mui/material'
import { FormProvider, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ILocationDTO, ILocationForm } from '../../../types/location'
import { NullDate } from '@/constants/defaultValues'
import { useAppDispatch } from '@/hooks/redux'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useCreateLocationMutation, useGetLastLocationQuery } from '../../../locationsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { Inputs } from './Inputs'

const def: ILocationForm = {
	needConfirm: false,
	department: '',
	person: '',
	dateOfIssue: dayjs().startOf('d').toISOString(),
}

type Props = {
	ids: string[]
}

export const Create: FC<Props> = ({ ids }) => {
	const [active, setActive] = useState(0)
	const savedIds = useRef(new Set<string>())
	const { palette } = useTheme()

	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetInstrumentByIdQuery(ids?.length ? ids[active] : '', {
		skip: !ids?.length || !ids[active],
	})
	const { data: loc, isFetching: isLocFetching } = useGetLastLocationQuery(data?.data.id || '', {
		skip: !data?.data.id,
	})
	const [create, { isLoading }] = useCreateLocationMutation()

	const methods = useForm<ILocationForm>({ defaultValues: def })

	useEffect(() => {
		if (loc?.data?.status == 'used') methods.setValue('isToReserve', true)
		else methods.setValue('isToReserve', false)
	}, [loc?.data?.status, methods])

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'NewLocation', isOpen: false }))
	}

	const activeHandler = (type: 'prev' | 'next') => () => {
		if (type == 'prev') setActive(prev => prev - 1)
		else setActive(prev => prev + 1)
	}

	const saveHandler = methods.handleSubmit(async form => {
		if (!data) return
		console.log('save', form, methods.formState.dirtyFields)

		if (form.dateOfIssue < (loc?.data?.dateOfIssue || 0)) {
			toast.error('Текущая дата выдачи инструмента меньше предыдущей')
			return
		}

		const dto: ILocationDTO = {
			instrumentId: data.data.id,
			person: form.isToReserve ? '' : form.person,
			department: form.isToReserve ? '' : form.department,
			dateOfIssue: form.dateOfIssue,
			dateOfReceiving: !form.needConfirm || form.isToReserve ? form.dateOfIssue : NullDate,
			needConfirm: form.needConfirm,
			status: !form.needConfirm ? (form.isToReserve ? 'reserve' : 'used') : 'moved',
		}

		try {
			await create(dto).unwrap()
			toast.success('Данные о перемещении добавлены')
			savedIds.current.add(ids[active])
			setActive(prev => prev + 1)
			if (savedIds.current.size == ids.length) closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	if (!ids?.length) return <Typography textAlign={'center'}>Инструменты не выбраны</Typography>
	return (
		<Stack position={'relative'} mt={-2.5}>
			{isFetching || isLoading || isLocFetching ? <BoxFallback /> : null}

			<Stack spacing={2} direction={'row'} paddingX={3} width={'100%'}>
				{ids.length > 1 && (
					<IconButton onClick={activeHandler('prev')} disabled={active == 0}>
						<LeftArrowIcon fontSize={16} fill={active == 0 ? palette.action.disabled : palette.grey[900]} />
					</IconButton>
				)}

				<Typography textAlign={'center'} fontSize={'1.2rem'} fontWeight={'bold'} sx={{ width: '100%' }}>
					<Typography component={'span'} mr={2.5} fontSize={'1.4rem'} color='primary'>
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
				{loc?.data.status != 'moved' ? (
					<FormProvider {...methods}>
						{loc?.data.status == 'used' ? <SendToReserve /> : <Inputs minDate={loc?.data.dateOfIssue} />}
					</FormProvider>
				) : (
					<Typography textAlign={'center'}>Подтверждение о получении еще не получено</Typography>
				)}

				<Divider sx={{ width: '50%', alignSelf: 'center', mt: 3 }} />
				<Stack spacing={2} direction={'row'} mt={2}>
					<Button onClick={closeHandler} variant='outlined' fullWidth>
						Отмена
					</Button>
					<Button
						disabled={savedIds.current.has(ids[active]) || loc?.data.status == 'moved'}
						type='submit'
						variant='contained'
						fullWidth
					>
						{loc?.data.status == 'used' ? 'Да' : 'Сохранить'}
					</Button>
				</Stack>
			</Stack>
		</Stack>
	)
}

const SendToReserve = () => {
	return <Typography>Переместить инструмент в резерв?</Typography>
}
