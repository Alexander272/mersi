import { FC, useEffect } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { Button, FormControl, InputLabel, MenuItem, Select, Stack, TextField, useTheme } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { ISection } from '../../types/sections'
import {
	useCreateSectionMutation,
	useDeleteSectionMutation,
	useGetGroupedSectionsQuery,
	useUpdateSectionMutation,
} from '../../sectionsApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'
import { useGetRealmsQuery } from '@/features/realms/realmsApiSlice'

type Props = {
	section: string
	setSection: (section: string) => void
}

type Form = Omit<ISection, 'id' | 'created'>

// TODO
const defaultValues: Form = {
	name: '',
	realmId: '',
	position: 1,
	bidType: '',
	verificationDay: 5,
}

export const Form: FC<Props> = ({ section, setSection }) => {
	const { palette } = useTheme()

	const { data: realms } = useGetRealmsQuery({ all: true })
	const { data, isFetching } = useGetGroupedSectionsQuery(null)

	const [create, { isLoading: isCreating }] = useCreateSectionMutation()
	const [update, { isLoading: isUpdating }] = useUpdateSectionMutation()
	const [remove, { isLoading: isDeleting }] = useDeleteSectionMutation()

	const {
		control,
		reset,
		handleSubmit,
		watch,
		formState: { dirtyFields },
	} = useForm<Form>({ values: defaultValues })
	const realmId = watch('realmId')

	useEffect(() => {
		if (data && section != 'new') {
			let selected: ISection | null = null
			data.data.forEach(item => {
				const tmp = item.sections.find(s => s.id == section)
				if (tmp) selected = tmp
			})
			if (selected) reset(selected)
		} else reset(defaultValues)
	}, [data, reset, section])

	const saveHandler = handleSubmit(async form => {
		if (!Object.keys(dirtyFields).length) return

		if (!data) {
			form.position = 1
		} else {
			const lastIndex = data?.data.length - 1
			const position = data?.data[lastIndex]?.sections[data?.data[lastIndex]?.sections.length - 1]?.position ?? 0
			if (section == 'new') form.position = position + 1
		}

		const newData = {
			...form,
			id: section || '',
			position: +form.position,
			maxPosition: +form.position,
			verificationDay: +form.verificationDay,
		}
		try {
			if (section == 'new') {
				const payload = await create(newData).unwrap()
				setSection(payload.id)
				toast.success('Область создана')
			} else {
				await update(newData).unwrap()
				toast.success('Область обновлена')
			}
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = async () => {
		if (section == 'new') return

		try {
			await remove(section).unwrap()
			setSection('new')
			toast.success('Область удалена')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	let length = data?.data.find(item => item.id == realmId)?.sections.length || 0
	if (section == 'new') length++

	return (
		<Stack component={'form'} onSubmit={saveHandler} position={'relative'}>
			{isFetching || isDeleting || isUpdating || isCreating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'name'}
					render={({ field }) => <TextField {...field} label={'Название'} fullWidth />}
				/>

				<FormControl fullWidth>
					<InputLabel id='realmId'>Область</InputLabel>
					<Controller
						control={control}
						name='realmId'
						render={({ field }) => (
							<Select labelId='realmId' {...field} label={'Область'}>
								<MenuItem value='' disabled>
									Выберите область
								</MenuItem>

								{realms?.data.map(item => (
									<MenuItem key={item.id} value={item.id}>
										{item.name}
									</MenuItem>
								))}
							</Select>
						)}
					/>
				</FormControl>
			</Stack>
			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'bidType'}
					render={({ field }) => <TextField {...field} label={'Тип'} fullWidth />}
				/>

				<Controller
					control={control}
					name={'verificationDay'}
					render={({ field }) => (
						<TextField {...field} label={'День месяца для уведомлений о поверке'} fullWidth />
					)}
				/>
			</Stack>

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2} alignItems={'center'}>
				<Controller
					control={control}
					name={'position'}
					render={({ field }) => (
						<TextField
							{...field}
							sx={{ maxWidth: '50%' }}
							label={'Позиция'}
							fullWidth
							slotProps={{
								htmlInput: {
									step: 1,
									min: 1,
									max: length,
									type: 'number',
								},
							}}
						/>
					)}
				/>

				<Stack
					direction={'row'}
					flexGrow={1}
					spacing={2}
					mb={1}
					justifyContent={'flex-end'}
					height={38}
					width={'100%'}
					sx={{ maxWidth: '50%' }}
				>
					<Button
						variant='outlined'
						type='submit'
						disabled={!Object.keys(dirtyFields).length}
						sx={{ minWidth: 56 }}
					>
						<SaveIcon
							fontSize={18}
							fill={!Object.keys(dirtyFields).length ? palette.action.disabled : palette.primary.main}
						/>
					</Button>

					<Confirm
						onClick={deleteHandler}
						buttonComponent={
							<Button
								variant='outlined'
								color='error'
								disabled={section == 'new'}
								sx={{ minWidth: 56, height: '100%' }}
							>
								<FileDeleteIcon
									fontSize={20}
									fill={section == 'new' ? palette.action.disabled : palette.error.main}
								/>
							</Button>
						}
						confirmText='Вы уверены, что хотите удалить секцию?'
						width='56'
					/>
				</Stack>
			</Stack>
		</Stack>
	)
}
