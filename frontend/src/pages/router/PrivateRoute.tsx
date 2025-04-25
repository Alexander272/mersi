import { Navigate, Outlet, useLocation } from 'react-router'

import { useAppSelector } from '@/hooks/redux'
import { getMenu, getToken } from '@/features/user/userSlice'
import { Forbidden } from '../forbidden/ForbiddenLazy'
import { AppRoutes } from './routes'

// проверка авторизации пользователя
export default function PrivateRoute() {
	const token = useAppSelector(getToken)
	const menu = useAppSelector(getMenu)
	const location = useLocation()

	if (!token) return <Navigate to={AppRoutes.Auth} state={{ from: location }} />
	if (!menu || !menu.length) return <Forbidden />

	return <Outlet />
}
