import { FC } from 'react'
import { Button, Divider, FormControl, Stack, TextField, Typography } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IChangePositionDTO } from '../../types/si'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'
import { useChangePositionMutation, useGetSIByIdQuery } from '../../siApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { LongRightIcon } from '@/components/Icons/LongRightIcon'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { useGetSI } from '../../hooks/getSI'

type Props = {
	id: string
}

export const ChangePositionForm: FC<Props> = ({ id }) => {
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetSIByIdQuery(id, { skip: !id })
	const { data: si } = useGetSI()
	const [change, { isLoading }] = useChangePositionMutation()

	const methods = useForm<IChangePositionDTO>({
		values: { sectionId: section?.id || '', oldPosition: data?.data.instrument.position || 1, newPosition: 1 },
	})

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'ChangePosition', isOpen: false }))
	}

	const saveHandler = methods.handleSubmit(async form => {
		console.log('save', form, methods.formState.dirtyFields)

		form.newPosition = +form.newPosition
		form.oldPosition = +form.oldPosition

		try {
			await change(form).unwrap()
			toast.success('Данные обновлены')
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack position={'relative'} mt={-5}>
			{isFetching || isLoading ? <BoxFallback /> : null}

			<Stack mt={2} component={'form'} onSubmit={saveHandler}>
				<Typography textAlign={'center'} mb={1.5}>
					{data?.data.instrument.name}
				</Typography>

				<Stack direction={'row'} spacing={2} alignItems={'center'} mb={2}>
					<FormControl fullWidth>
						<Typography textAlign={'center'} mb={0.5}>
							Текущая позиция
						</Typography>
						<Controller
							control={methods.control}
							name={'oldPosition'}
							rules={{ required: true }}
							render={({ field, fieldState: { error } }) => (
								<TextField
									{...field}
									value={field.value || ''}
									onChange={e => field.onChange(+(e.target.value || 0))}
									fullWidth
									disabled
									error={Boolean(error)}
									slotProps={{
										htmlInput: {
											type: 'number',
											step: 1,
											min: 1,
											max: si?.total,
										},
									}}
								/>
							)}
						/>
					</FormControl>

					<LongRightIcon fontSize={28} width={30} height={44} alignSelf={'flex-end'} />

					<FormControl fullWidth>
						<Typography textAlign={'center'} mb={0.5}>
							Новая позиция
						</Typography>
						<Controller
							control={methods.control}
							name={'newPosition'}
							rules={{ required: true }}
							render={({ field, fieldState: { error } }) => (
								<TextField
									{...field}
									value={field.value || ''}
									onChange={e => field.onChange(+(e.target.value || 0))}
									fullWidth
									error={Boolean(error)}
									slotProps={{
										htmlInput: {
											type: 'number',
											step: 1,
											min: 1,
											max: si?.total,
										},
									}}
								/>
							)}
						/>
					</FormControl>
				</Stack>

				<Divider sx={{ width: '50%', alignSelf: 'center' }} />
				<Stack spacing={2} direction={'row'} mt={2}>
					<Button onClick={closeHandler} variant='outlined' fullWidth>
						Отмена
					</Button>
					<Button type='submit' variant='contained' fullWidth>
						Применить
					</Button>
				</Stack>
			</Stack>
		</Stack>
	)
}
