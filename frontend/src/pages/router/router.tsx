import { createBrowserRouter, RouteObject } from 'react-router'

import { Layout } from '@/components/Layout/Layout'
import { NotFound } from '@/pages/notFound/NotFoundLazy'
import { Auth } from '@/pages/auth/AuthLazy'
import { Home } from '@/pages/home/HomeLazy'
import { Sections } from '@/pages/sections/SectionLazy'
import { AppRoutes } from './routes'
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
					// {
					// 	path: AppRoutes.EMPLOYEES,
					// 	element: <Employees />,
					// },
					// {
					// 	path: AppRoutes.PLACES,
					// 	element: <Places />,
					// },
					// {
					// 	path: AppRoutes.REALMS,
					// 	element: <Realms />,
					// },
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
