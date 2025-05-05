import { Tab, Tabs } from '@mui/material'
import { FC, SyntheticEvent } from 'react'

type Props = {
	value: string
	onChange: (value: string) => void
}

export const ChildrenTabs: FC<Props> = ({ value, onChange }) => {
	const tabHandler = (_event: SyntheticEvent, value: string) => {
		onChange(value)
	}

	return (
		<Tabs
			value={value}
			onChange={tabHandler}
			variant='scrollable'
			sx={{
				mb: 2,
				'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
				'.MuiTabs-scrollButtons.Mui-disabled': {
					height: 0,
				},
			}}
		>
			<Tab
				label='Колонки'
				value='columns'
				sx={{
					textTransform: 'inherit',
					borderRadius: 3,
					transition: 'all 0.3s ease-in-out',
					maxWidth: '50%',
					minHeight: 48,
					flexGrow: 1,
					':hover': {
						backgroundColor: '#f5f5f5',
					},
				}}
			/>
			<Tab
				label='Форма'
				value='form'
				sx={{
					textTransform: 'inherit',
					borderRadius: 3,
					transition: 'all 0.3s ease-in-out',
					maxWidth: '50%',
					minHeight: 48,
					flexGrow: 1,
					':hover': {
						backgroundColor: '#f5f5f5',
					},
				}}
			/>
		</Tabs>
	)
}
