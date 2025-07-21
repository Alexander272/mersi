import { FC } from 'react'
import { Box } from '@mui/material'

import { useGetRepairQuery } from '../../repair/repairApiSlice'
import { Formatter } from '@/features/table/utils/formatter'
import { NoRowsOverlay } from '@/features/table/components/NoRowsOverlay/components/NoRowsOverlay'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Table } from '@/components/Table/Table'
import { TableHead } from '@/components/Table/TableHead'
import { TableRow } from '@/components/Table/TableRow'
import { TableBody } from '@/components/Table/TableBody'
import { TableCell } from '@/components/Table/TableCell'

type Props = {
	instrumentId: string
}

export const Repair: FC<Props> = ({ instrumentId }) => {
	const { data, isFetching } = useGetRepairQuery(instrumentId, { skip: !instrumentId })

	if (isFetching) return <BoxFallback />
	return (
		<Table>
			<TableHead>
				<TableRow height={40}>
					<TableCell width={288}>Дефект</TableCell>
					<TableCell width={288}>Ремонт</TableCell>
					<TableCell width={288}>Период ремонта</TableCell>
					<TableCell width={288}>Описание</TableCell>
				</TableRow>
			</TableHead>
			<Box maxHeight={350} overflow={'auto'} position={'relative'} minHeight={150}>
				{!data?.data.length && <NoRowsOverlay />}
				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id} height={38}>
							<TableCell width={288}>{Formatter('text', item.defect)}</TableCell>
							<TableCell width={288}>{Formatter('text', item.work)}</TableCell>
							<TableCell width={288}>{`${Formatter('date', item.periodStart)} - ${Formatter(
								'date',
								item.periodEnd
							)}`}</TableCell>
							<TableCell width={288}>{Formatter('text', item.description)}</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Box>
		</Table>
	)
}
