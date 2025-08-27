import { FC } from 'react'
import { Box } from '@mui/material'
import dayjs from 'dayjs'

import { DayjsFormat } from '@/constants/dateFormat'
import { useGetLocationsQuery } from '../../locations/locationsApiSlice'
import { NoRowsOverlay } from '@/features/table/components/NoRowsOverlay/components/NoRowsOverlay'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Table } from '@/components/Table/Table'
import { TableHead } from '@/components/Table/TableHead'
import { TableBody } from '@/components/Table/TableBody'
import { TableRow } from '@/components/Table/TableRow'
import { TableCell } from '@/components/Table/TableCell'
import { CellText } from '@/components/CellText/CellText'

type Props = {
	instrumentId: string
}

export const Location: FC<Props> = ({ instrumentId }) => {
	const { data, isFetching } = useGetLocationsQuery(instrumentId, { skip: !instrumentId })

	if (isFetching) return <BoxFallback />
	return (
		<Table>
			<TableHead>
				<TableRow height={70}>
					<TableCell width={230}>Дата выдачи</TableCell>
					<TableCell width={230}>Дата получения</TableCell>
					<TableCell width={230}>Место нахождения</TableCell>
					<TableCell width={230}>Подтверждение запрашивалось?</TableCell>
					<TableCell width={230}>Подтверждение получено?</TableCell>
				</TableRow>
			</TableHead>

			<Box maxHeight={350} overflow={'auto'} position={'relative'} minHeight={150}>
				{!data?.data.length && <NoRowsOverlay />}

				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id} sx={{ minHeight: 38, cursor: 'default' }}>
							<TableCell width={230}>
								<CellText value={dayjs(item.dateOfIssue * 1000).format(DayjsFormat)} />
							</TableCell>
							<TableCell width={230}>
								<CellText
									value={
										item.dateOfReceiving
											? dayjs(item.dateOfReceiving * 1000).format(DayjsFormat)
											: '-'
									}
								/>
							</TableCell>
							<TableCell width={230}>
								<CellText
									value={item.status == 'reserve' && !item.place ? 'Резерв' : item.place || ''}
								/>
							</TableCell>
							<TableCell width={230}>
								<CellText value={item.needConfirm ? 'Да' : 'Нет'} />
							</TableCell>
							<TableCell width={230}>
								<CellText value={item.hasConfirmed ? 'Да' : 'Нет'} />
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Box>
		</Table>
	)
}
