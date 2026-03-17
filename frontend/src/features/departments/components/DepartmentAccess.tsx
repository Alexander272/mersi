import { FC, useMemo } from 'react'
import { Autocomplete, Box, Button, Stack, TextField, Tooltip, Typography, useTheme } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IChangeDepAccess, IDepartment } from '../types/departments'
import { useGetUserByAccessQuery } from '@/features/user/usersApiSlice'
import { useChangeDepAccessesMutation, useGetDepAccessesQuery } from '../departmentAccessApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { QuestionIcon } from '@/components/Icons/QuestionIcon'
import { SaveIcon } from '@/components/Icons/SaveIcon'

type Props = {
	department?: IDepartment
}

export const DepartmentAccess: FC<Props> = ({ department }) => {
	const { palette } = useTheme()
	const depId = department?.id || ''

	const { data: accessData, isFetching: isAccessLoading } = useGetDepAccessesQuery(
		{ department: depId },
		{ skip: !department },
	)
	const { data: usersData, isFetching: isUsersLoading } = useGetUserByAccessQuery(null)
	const [change, { isLoading: isUpdating }] = useChangeDepAccessesMutation()

	const userOptions = useMemo(() => usersData?.data || [], [usersData])

	const {
		control,
		handleSubmit,
		formState: { isDirty },
		reset,
	} = useForm<IChangeDepAccess>({
		values: useMemo(
			() => ({
				departmentId: depId,
				userIds: accessData?.data?.map(el => el.userId) || [],
			}),
			[accessData, depId],
		),
	})

	const saveHandler = handleSubmit(async form => {
		console.log('save access', form)
		try {
			await change(form).unwrap()
			reset(form)
			toast.success('Доступы сохранены')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const isLoading = isAccessLoading || isUsersLoading || isUpdating
	return (
		<Stack direction='row' position='relative' width='100%'>
			{isLoading && <Fallback position='absolute' zIndex={5} background='#f5f5f557' />}

			<Stack
				spacing={2}
				mt={1}
				mr={1}
				direction='row'
				component='form'
				onSubmit={saveHandler}
				width='100%'
				alignItems='center'
			>
				<Stack direction='row' alignItems='center' spacing={1}>
					<Typography whiteSpace='nowrap'>Доступ: </Typography>
					<Tooltip title='Выбор пользователей, которые могут видеть инструменты в данном подразделении'>
						<Box
							sx={{
								cursor: 'help',
								width: 20,
								height: 20,
								display: 'flex',
								borderRadius: '50%',
								background: '#eff3fc',
								justifyContent: 'center',
								alignItems: 'center',
							}}
						>
							<QuestionIcon fontSize={12} />
						</Box>
					</Tooltip>
				</Stack>

				<Controller
					name='userIds'
					control={control}
					render={({ field: { onChange, value } }) => (
						<Autocomplete
							multiple
							options={userOptions}
							getOptionLabel={option => `${option.lastName} ${option.firstName}`}
							value={userOptions.filter(user => value?.includes(user.id))}
							onChange={(_, newValue) => {
								onChange(newValue.map(user => user.id))
							}}
							renderInput={params => (
								<TextField {...params} label='Пользователи' variant='outlined' size='small' />
							)}
							fullWidth
							disableCloseOnSelect
							limitTags={3}
						/>
					)}
				/>

				<Button
					variant='outlined'
					type='submit'
					disabled={!isDirty || isLoading}
					sx={{ minWidth: 56, height: 40 }}
				>
					<SaveIcon fontSize={18} fill={!isDirty ? palette.action.disabled : palette.primary.main} />
				</Button>
			</Stack>
		</Stack>
	)
}
