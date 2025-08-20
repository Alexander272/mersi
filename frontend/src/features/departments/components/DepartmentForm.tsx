import { FC, useEffect } from 'react'
import { Button, FormControl, InputLabel, MenuItem, Select, Stack, TextField, useTheme } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IDepartment } from '../types/departments'
import { useGetAllChannelsQuery } from '@/features/channel/channelApiSlice'
import {
	useCreateDepartmentMutation,
	useGetDepartmentsByIdQuery,
	useDeleteDepartmentMutation,
	useUpdateDepartmentMutation,
} from '../departmentApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { Confirm } from '@/components/Confirm/Confirm'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { FileDeleteIcon } from '@/components/Icons/FileDeleteIcon'

type Props = {
	department: string
	setDepartment: (department: string) => void
}

const defaultValues = { name: '', channelId: '' }

export const DepartmentForm: FC<Props> = ({ department, setDepartment }) => {
	const { data: channels, isFetching: channelsIsFetching } = useGetAllChannelsQuery(null)
	const { data, isFetching } = useGetDepartmentsByIdQuery(department, { skip: !department || department == 'new' })

	const [create] = useCreateDepartmentMutation()
	const [update] = useUpdateDepartmentMutation()
	const [remove] = useDeleteDepartmentMutation()

	const { palette } = useTheme()

	const {
		control,
		reset,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<Pick<IDepartment, 'name' | 'channelId'>>({ values: defaultValues })

	useEffect(() => {
		if (data && department != 'new') reset({ name: data.data.name, channelId: data.data.channelId })
		else reset(defaultValues)
	}, [data, channels, reset, department])

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)
		if (!Object.keys(dirtyFields).length) return

		const newData = { ...form, id: data?.data.id || '', leaderId: data?.data.leaderId || '', channelName: '' }
		try {
			if (department == 'new') {
				const payload = await create(newData).unwrap()
				setDepartment(payload.id)
				toast.success('Подразделение создано')
			} else {
				await update(newData).unwrap()
				toast.success('Подразделение обновлено')
			}
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const deleteHandler = async () => {
		if (department == 'new') return
		try {
			await remove(department).unwrap()
			setDepartment('new')
			toast.success('Подразделение удалено')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	// if (isFetching || channelsIsFetching) return <Fallback />
	return (
		<Stack component={'form'} onSubmit={saveHandler} alignItems={'center'} position={'relative'}>
			{isFetching || channelsIsFetching ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />
			) : null}

			<Stack direction={'row'} flexGrow={1} width={'100%'} spacing={2} mb={2}>
				<Controller
					control={control}
					name={'name'}
					render={({ field }) => <TextField {...field} label={'Подразделение'} fullWidth />}
				/>

				<Button
					variant='outlined'
					type='submit'
					disabled={!Object.keys(dirtyFields).length}
					sx={{
						minWidth: 56,
					}}
				>
					<SaveIcon
						fontSize={18}
						fill={!Object.keys(dirtyFields).length ? palette.action.disabled : palette.primary.main}
					/>
				</Button>
			</Stack>

			<Stack direction={'row'} flexGrow={1} width={'100%'} spacing={2}>
				<FormControl fullWidth>
					<InputLabel id='channel'>Канал для уведомлений</InputLabel>
					<Controller
						control={control}
						name={'channelId'}
						render={({ field }) => (
							<Select {...field} labelId='channel' label='Канал для уведомлений'>
								{channels?.data.map(c => (
									<MenuItem value={c.id} key={c.id}>
										{c.name}
									</MenuItem>
								))}
							</Select>
						)}
					/>
				</FormControl>
				{/* //TODO добавить подсказку для чего это поле */}

				{/* {department != 'new' && ( */}
				<Confirm
					onClick={deleteHandler}
					buttonComponent={
						<Button
							variant='outlined'
							color='error'
							disabled={department == 'new'}
							sx={{ minWidth: 56, height: '100%' }}
						>
							<FileDeleteIcon
								fontSize={20}
								fill={department == 'new' ? palette.action.disabled : palette.error.main}
							/>
						</Button>
					}
					confirmText='Вы уверены, что хотите удалить подразделение?'
					width='56'
				/>
				{/* )} */}
			</Stack>
		</Stack>
	)
}
