import { FC, SyntheticEvent } from 'react'
import { Stack, Tab, Tabs, useTheme } from '@mui/material'

import { Fallback } from '@/components/Fallback/Fallback'
import { useGetRealmsQuery } from '../realmsApiSlice'

type Props = {
	realm: string
	setRealm: (realm: string) => void
}

export const RealmsList: FC<Props> = ({ realm, setRealm }) => {
	const { palette } = useTheme()

	const { data, isFetching } = useGetRealmsQuery({ all: true })

	const tabHandler = (_event: SyntheticEvent, value: string) => {
		setRealm(value)
	}

	return (
		<Stack position={'relative'} minWidth={280} height={'100%'}>
			{isFetching && <Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />}
			<Tabs
				orientation='vertical'
				value={realm}
				onChange={tabHandler}
				variant='scrollable'
				sx={{
					borderRight: 1,
					borderColor: 'divider',
					maxHeight: 760,
					height: '100%',
					'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
					'.MuiTabs-scrollButtons.Mui-disabled': {
						height: 0,
					},
				}}
			>
				<Tab
					label='Добавить'
					value='new'
					sx={{
						mt: 0.5,
						textTransform: 'inherit',
						borderRadius: 3,
						transition: 'all 0.3s ease-in-out',
						maxWidth: '100%',
						minHeight: 44,
						backgroundColor: '#9ab2ef29',
						color: palette.primary.main,
						':hover': {
							backgroundColor: '#9ab2ef58',
						},
					}}
				/>

				{data?.data.map(m => (
					<Tab
						key={m.id}
						label={m.name}
						value={m.id}
						sx={{
							textTransform: 'inherit',
							borderRadius: 3,
							transition: 'all 0.3s ease-in-out',
							maxWidth: '100%',
							minHeight: 48,
							':hover': {
								backgroundColor: '#f5f5f5',
							},
						}}
					/>
				))}
			</Tabs>
		</Stack>
	)
}
