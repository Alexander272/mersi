import { AppBar, Box, Stack, Toolbar } from '@mui/material'

import { PermRules } from '@/constants/permissions'
import { AppRoutes } from '@/pages/router/routes'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useSignOutMutation } from '@/features/auth/authApiSlice'
import { getToken } from '@/features/user/userSlice'
import { ActiveRealm } from '@/features/realms/components/ActiveRealm'
import { GeometryIcon } from '../Icons/GeometryIcon'
import { NavButton, NavLink } from './header.style'

import logo from '@/assets/logo.webp'

export const LayoutHeader = () => {
	const [signOut] = useSignOutMutation()

	const token = useAppSelector(getToken)

	const showRealmsSetting = useCheckPermission(PermRules.Realms.Write)

	const signOutHandler = () => {
		void signOut(null)
	}

	return (
		<AppBar sx={{ borderRadius: 0 }}>
			<Toolbar sx={{ justifyContent: 'space-between', alignItems: 'inherit' }}>
				<Box alignSelf={'center'} display={'flex'} alignItems={'center'}>
					<img height={46} width={157} src={logo} alt='logo' />
					<GeometryIcon fill={'#042245'} />
				</Box>

				{token && (
					<Stack direction={'row'} spacing={3} minHeight={'100%'}>
						<NavButton>
							<ActiveRealm sx={{ mr: -0.5 }} />
						</NavButton>
						{showRealmsSetting && <NavLink to={AppRoutes.Realm}>Области</NavLink>}
						<NavButton onClick={signOutHandler}>Выйти</NavButton>
					</Stack>
				)}
			</Toolbar>
		</AppBar>
	)
}
