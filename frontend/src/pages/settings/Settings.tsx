import { Suspense, SyntheticEvent, useState } from 'react'
import { Box, Breadcrumbs, Tab, Tabs } from '@mui/material'

import { AppRoutes } from '../router/routes'
import { PageBox } from '@/components/PageBox/PageBox'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { Fallback } from '@/components/Fallback/Fallback'
import { Realms } from './realms/RealmsLazy'
import { Sections } from './sections/SectionLazy'
import { Statuses } from './statuses/StatusesLazy'
import { ContextMenu } from './contextMenu/ContextMenuLazy'
import { Tools } from './tools/ToolsLazy'
import { Verification } from './verification/VerificationLazy'
import History from './history/History'

const tabsData = [
	{ id: 'realm', name: 'Области' },
	{ id: 'sections', name: 'Секции' },
	{ id: 'status', name: 'Статусы' },
	{ id: 'context', name: 'Контекстное меню' },
	{ id: 'tools', name: 'Меню "Инструменты"' },
	{ id: 'verification', name: 'Поверки' },
	{ id: 'history', name: 'История' },
]

export default function Settings() {
	const [active, setActive] = useState('realm')

	const tabHandler = (_event: SyntheticEvent, value: string) => {
		setActive(value)
	}

	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingY={2}
				margin={'0 auto'}
				width={{ xl: '66%', lg: '86%', md: '100%' }}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				// flexGrow={1}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: 2, px: 2 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.Settings} active>
						Настройки
					</Breadcrumb>
				</Breadcrumbs>

				<Tabs
					value={active}
					onChange={tabHandler}
					// variant='scrollable'
					// scrollButtons='auto'
					sx={{
						borderBottom: 1,
						borderColor: 'divider',
						width: '100%',
						mb: 3,
						px: 2,
						'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
						'.MuiTabs-scrollButtons.Mui-disabled': {
							height: 0,
						},
					}}
				>
					{tabsData.map(t => (
						<Tab
							key={t.id}
							label={t.name}
							value={t.id}
							sx={{
								textTransform: 'inherit',
								borderRadius: 3,
								transition: 'all 0.3s ease-in-out',
								maxWidth: '100%',
								flexGrow: 1,
								// minHeight: 48,
								':hover': {
									backgroundColor: '#f5f5f5',
								},
							}}
						/>
					))}
				</Tabs>

				<Suspense fallback={<Fallback />}>
					{active == 'realm' && <Realms />}
					{active == 'sections' && <Sections />}
					{active == 'status' && <Statuses />}
					{active == 'context' && <ContextMenu />}
					{active == 'tools' && <Tools />}
					{active == 'verification' && <Verification />}
					{active == 'history' && <History />}
				</Suspense>
			</Box>
		</PageBox>
	)
}
