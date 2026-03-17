import { SyntheticEvent, useState } from 'react'
import { Stack, Tab, Tabs } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { VerificationFieldsList } from '@/features/table/modules/verification/modules/form/List'
import { VerificationHistoryList } from '@/features/table/modules/verification/modules/history/List'

const tabsData = [
	{ id: 'form', name: 'Форма' },
	{ id: 'history', name: 'История' },
]

export default function Verification() {
	const [item, setItem] = useState<string>('new')
	const [active, setActive] = useState('form')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	const tabHandler = (_event: SyntheticEvent, value: string) => {
		setActive(value)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'} pl={2}>
			<List item={item} setItem={itemHandler} readonly />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1, pr: 2 }}>
				<Tabs
					value={active}
					onChange={tabHandler}
					sx={{
						borderBottom: 1,
						borderColor: 'divider',
						width: '100%',
						mb: 3,
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
								':hover': {
									backgroundColor: '#f5f5f5',
								},
							}}
						/>
					))}
				</Tabs>

				{active == 'form' && <VerificationFieldsList section={item} />}
				{active == 'history' && <VerificationHistoryList section={item} />}
			</Stack>
		</Stack>
	)
}
