import { useState } from 'react'
import { Button, InputAdornment, LinearProgress, Stack, TextField, Typography, useTheme } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { ISignIn } from '../types/auth'
import { useAppDispatch } from '@/hooks/redux'
import { setUser } from '@/features/user/userSlice'
import { VisibleIcon } from '@/components/Icons/VisibleIcon'
import { InVisibleIcon } from '@/components/Icons/InVisibleIcon'
import { useSignInMutation } from '../authApiSlice'

const defaultValues: ISignIn = { username: '', password: '' }

export const SignInForm = () => {
	const [passIsVisible, setPassIsVisible] = useState(false)
	const { palette } = useTheme()

	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { errors },
	} = useForm<ISignIn>({ defaultValues })

	const [signIn, { isLoading }] = useSignInMutation()

	const togglePassVisible = () => setPassIsVisible(prev => !prev)

	const signInHandler = async (data: ISignIn) => {
		try {
			const payload = await signIn(data).unwrap()
			dispatch(setUser(payload.data))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<Stack component='form' onSubmit={handleSubmit(signInHandler)} sx={{ position: 'relative' }}>
			{isLoading ? <LinearProgress sx={{ position: 'absolute', bottom: -20, left: 0, right: 0 }} /> : null}

			<Typography
				variant='h2'
				align='center'
				fontSize={'1.5rem'}
				color={palette.primary.main}
				paddingBottom={1.25}
				mb={1.25}
				fontWeight={'bold'}
				lineHeight={'inherit'}
				sx={{ borderBottom: '1px solid #e5e4e9', letterSpacing: '1.2px' }}
			>
				Вход
			</Typography>

			<Stack spacing={2} marginY={2}>
				<Controller
					control={control}
					name='username'
					rules={{ required: true }}
					render={({ field }) => (
						<TextField
							name={field.name}
							value={field.value}
							onChange={field.onChange}
							placeholder='Имя пользователя'
							error={Boolean(errors.username)}
							helperText={errors.username ? 'Поле не может быть пустым' : ''}
							disabled={isLoading}
							sx={{ '& .MuiOutlinedInput-root': { borderRadius: 10 } }}
						/>
					)}
				/>

				<Controller
					control={control}
					name='password'
					rules={{ required: true }}
					render={({ field }) => (
						<TextField
							name={field.name}
							value={field.value}
							onChange={field.onChange}
							type={passIsVisible ? 'text' : 'password'}
							placeholder='Пароль'
							error={Boolean(errors.password)}
							helperText={errors.password ? 'Поле не может быть пустым' : ''}
							disabled={isLoading}
							sx={{ '& .MuiOutlinedInput-root': { borderRadius: 10, paddingRight: 0.5 } }}
							InputProps={{
								endAdornment: (
									<InputAdornment
										position='start'
										onClick={togglePassVisible}
										sx={{ cursor: 'pointer' }}
									>
										{passIsVisible ? <VisibleIcon /> : <InVisibleIcon />}
									</InputAdornment>
								),
							}}
						/>
					)}
				/>
			</Stack>

			<Button type='submit' disabled={isLoading} variant='contained' sx={{ borderRadius: 10, marginY: 3 }}>
				Войти
			</Button>
		</Stack>
	)
}
