import { FC, useEffect } from 'react'
import { Controller, useForm } from 'react-hook-form'
import {
	Box,
	Button,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	Slider,
	Stack,
	TextField,
	Typography,
	useTheme,
} from '@mui/material'

import type { ISection } from '../../types/sections'
import { Fallback } from '@/components/Fallback/Fallback'
import { useGetGroupedSectionsQuery } from '../../sectionsApiSlice'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { Confirm } from '@/components/Confirm/Confirm'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'

type Props = {
	section: string
	setSection: (section: string) => void
}

type Form = Omit<ISection, 'id' | 'created'>

const defaultValues: Form = {
	name: '',
	realmId: '',
	position: 1,
}

export const Form: FC<Props> = ({ section, setSection }) => {
	const { palette } = useTheme()

	const { data, isFetching } = useGetGroupedSectionsQuery(null)

	const isCreating = false
	const isUpdating = false
	const isDeleting = false

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
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		//TODO дописать
	})

	const deleteHandler = async () => {
		if (section == 'new') return
		// try {
		// 	await remove(section).unwrap()
		setSection('new')
		// 	toast.success('Область удалена')
		// } catch (error) {
		// 	const fetchError = error as IFetchError
		// 	toast.error(fetchError.data.message, { autoClose: false })
		// }
	}

	let length = data?.data.find(item => item.id == realmId)?.sections.length || 0
	if (section == 'new') length++

	return (
		<Stack component={'form'} onSubmit={saveHandler} position={'relative'} mr={1}>
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

								{data?.data.map(item => (
									<MenuItem key={item.id} value={item.id}>
										{item.title}
									</MenuItem>
								))}
							</Select>
						)}
					/>
				</FormControl>
			</Stack>

			<Stack direction={'row'} flexGrow={1} spacing={2} mb={2} alignItems={'center'}>
				<Stack flexGrow={1} maxWidth={500}>
					<Typography>Позиция</Typography>
					<Stack direction={'row'} flexGrow={1}>
						<Box width={'90%'} paddingX={3}>
							<Controller
								control={control}
								name={'position'}
								render={({ field }) => (
									<Slider
										{...field}
										aria-label='position'
										defaultValue={length + 1}
										valueLabelDisplay='auto'
										shiftStep={1}
										step={1}
										marks
										min={1}
										max={length}
									/>
								)}
							/>
						</Box>

						<Controller
							control={control}
							name={'position'}
							render={({ field }) => (
								<TextField
									{...field}
									variant='standard'
									sx={{ width: 40 }}
									slotProps={{
										htmlInput: {
											step: 1,
											min: 1,
											max: length,
											type: 'number',
											'aria-labelledby': 'input-slider',
										},
									}}
								/>
							)}
						/>
					</Stack>
				</Stack>

				<Stack direction={'row'} flexGrow={1} spacing={2} mb={1} justifyContent={'flex-end'} height={38}>
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
