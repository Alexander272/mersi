import { FC, useEffect, useState } from 'react'
import { Stack, Tab, Tabs, Typography } from '@mui/material'

import { useAppSelector } from '@/hooks/redux'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { useGetHistoryTypesQuery } from '../historyApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Verification } from './Verification'
import { Repair } from './Repair'
import { Preservation } from './Preservation'
import { Save } from './Save'

type Props = {
	instrumentId: string
}

export const History: FC<Props> = ({ instrumentId }) => {
	const [value, setValue] = useState('')
	const section = useAppSelector(getSection)

	const { data: instrument, isFetching } = useGetInstrumentByIdQuery(instrumentId, { skip: !instrumentId })
	const { data, isFetching: isFetchingTypes } = useGetHistoryTypesQuery(section?.id || '', {
		skip: !section?.id,
	})

	useEffect(() => {
		if (data?.data) setValue(data.data[0].group)
	}, [data])

	const tabHandler = (_event: React.SyntheticEvent, newValue: string) => {
		setValue(newValue)
	}

	return (
		<Stack position={'relative'} mt={-2.5} minHeight={200}>
			{isFetching || isFetchingTypes ? <BoxFallback /> : null}

			<Typography textAlign={'center'} sx={{ width: '100%' }}>
				{instrument?.data.name} ({instrument?.data.factoryNumber})
			</Typography>

			{(data?.data.length || 0) > 1 && value != '' ? (
				<Tabs
					value={value}
					onChange={tabHandler}
					// variant='scrollable'
					sx={{
						mt: 1,
						mb: 2,
						borderBottom: 1,
						borderColor: '#00000014',
						'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
						'.MuiTabs-scrollButtons.Mui-disabled': {
							height: 0,
						},
					}}
				>
					{data?.data.map(f => (
						<Tab
							key={f.id}
							label={f.label}
							value={f.group}
							sx={{
								textTransform: 'inherit',
								borderRadius: 3,
								transition: 'all 0.3s ease-in-out',
								minHeight: 48,
								flexGrow: 1,
								':hover': {
									backgroundColor: '#f5f5f5',
								},
							}}
						/>
					))}
				</Tabs>
			) : null}

			{value == 'verifications' && <Verification instrumentId={instrumentId} />}
			{value == 'repair' && <Repair instrumentId={instrumentId} />}
			{value == 'preservation' && <Preservation instrumentId={instrumentId} />}
			{value == 'save' && <Save instrumentId={instrumentId} />}
		</Stack>
	)
}
