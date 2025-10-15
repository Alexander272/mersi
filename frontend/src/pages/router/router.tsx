import { createBrowserRouter, RouteObject } from 'react-router'

import { AppRoutes } from './routes'
import { Layout } from '@/components/Layout/Layout'
import { NotFound } from '@/pages/notFound/NotFoundLazy'
import { Auth } from '@/pages/auth/AuthLazy'
import { Home } from '@/pages/home/HomeLazy'
import { Sections } from '@/pages/sections/SectionLazy'
import { Realms } from '@/pages/realms/RealmsLazy'
import { Places } from '@/pages/places/PlacesLazy'
import { Employees } from '@/pages/employee/EmployeesLazy'
import { Import } from '@/pages/import/ImportLazy'
import { Settings } from '@/pages/settings/SettingsLazy'
import PrivateRoute from './PrivateRoute'

const config: RouteObject[] = [
	{
		element: <Layout />,
		errorElement: <NotFound />,
		children: [
			{
				path: AppRoutes.Auth,
				element: <Auth />,
			},
			{
				path: AppRoutes.Home,
				element: <PrivateRoute />,
				children: [
					{
						index: true,
						element: <Home />,
					},
					{
						path: AppRoutes.Employees,
						element: <Employees />,
					},
					{
						path: AppRoutes.Places,
						element: <Places />,
					},
					{
						path: AppRoutes.Import,
						element: <Import />,
					},
					{
						path: AppRoutes.Settings,
						element: <Settings />,
					},
					{
						path: AppRoutes.Realm,
						element: <Realms />,
					},
					{
						path: AppRoutes.Sections,
						element: <Sections />,
					},
				],
			},
		],
	},
]

export const router = createBrowserRouter(config)
