import { useState } from 'react'
import { Box, Breadcrumbs, Stack } from '@mui/material'

import { AppRoutes } from '../router/routes'
import { List } from '@/features/sections/components/List/List'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { PageBox } from '@/components/PageBox/PageBox'
import { Form } from '@/features/sections/components/Form/Form'
import { Columns } from '@/features/sections/modules/columns/components/Columns/Columns'
import { ChildrenTabs } from '@/features/sections/components/Tabs/ChildrenTabs'
import { FieldsList } from '@/features/sections/modules/form/components/List/FieldsList'

export default function Sections() {
	const [item, setItem] = useState('new')
	const [child, setChild] = useState('columns')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<PageBox>
			<Box
				borderRadius={3}
				padding={2}
				margin={'0 auto'}
				width={{ xl: '66%', lg: '86%', md: '100%' }}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				flexGrow={1}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: 2 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.Sections} active>
						Секции
					</Breadcrumb>
				</Breadcrumbs>

				<Stack direction={'row'} spacing={2} height={'100%'}>
					<List item={item} setItem={itemHandler} />
					<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
						<Form section={item} setSection={itemHandler} />
						{item != 'new' && (
							<>
								<ChildrenTabs value={child} onChange={setChild} />
								{child == 'columns' && <Columns section={item} />}
								{child == 'form' && <FieldsList section={item} />}
							</>
						)}
					</Stack>
				</Stack>
			</Box>
		</PageBox>
	)
}
