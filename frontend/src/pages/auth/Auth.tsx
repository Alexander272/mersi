import { useEffect } from 'react'
import { Box, Typography, useTheme } from '@mui/material'
import { useLocation, useNavigate } from 'react-router'

import { useAppSelector } from '@/hooks/redux'
import { SignInForm } from '@/features/auth/components/SignInForm'
import { getToken } from '@/features/user/userSlice'
import { GeometryIcon } from '@/components/Icons/GeometryIcon'
import { PageBox } from '@/components/PageBox/PageBox'

type LocationState = {
	from?: Location
}

export default function Auth() {
	const { palette } = useTheme()

	const navigate = useNavigate()
	const location = useLocation()

	const token = useAppSelector(getToken)

	useEffect(() => {
		const to: string = (location.state as LocationState)?.from?.pathname || '/'
		if (token) navigate(to, { replace: true })
	}, [token, navigate, location.state])

	return (
		<PageBox>
			<Box display={'flex'} flexDirection={'column'} justifyContent={'center'} alignItems={'center'} flexGrow={1}>
				<Typography
					mb={2}
					mt={-6}
					fontSize={'3rem'}
					display={'flex'}
					alignItems={'center'}
					gap={'8px'}
					color={'#042245'}
					fontFamily={'serif'}
				>
					МЕРСИ <GeometryIcon fontSize={'inherit'} fill={'#042245'} />
				</Typography>

				<Box
					marginX={3}
					borderRadius={4}
					paddingY={2.5}
					paddingX={3.75}
					width={{ sm: 400, xs: '100%' }}
					boxShadow={'2px 2px 8px 0px #3636362b'}
					sx={{ background: palette.background.paper }}
				>
					<SignInForm />
				</Box>
			</Box>
		</PageBox>
	)
}
