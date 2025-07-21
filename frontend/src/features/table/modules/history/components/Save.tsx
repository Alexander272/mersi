import { FC } from 'react'
import { Box } from '@mui/material'

import { useGetTransferToSaveQuery } from '../../transferToSave/transferApiSlice'
import { Formatter } from '@/features/table/utils/formatter'
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

export const Save: FC<Props> = ({ instrumentId }) => {
	const { data, isFetching } = useGetTransferToSaveQuery(instrumentId, { skip: !instrumentId })

	if (isFetching) return <BoxFallback />
	return (
		<Table>
			<TableHead>
				<TableRow height={40}>
					<TableCell width={300}>Передача / Возврат</TableCell>
					<TableCell width={850}>Примечание</TableCell>
				</TableRow>
			</TableHead>
			<Box maxHeight={350} overflow={'auto'} position={'relative'} minHeight={150}>
				{!data?.data.length && <NoRowsOverlay />}
				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id} height={38}>
							<TableCell width={300}>
								{Formatter('date', item.dateStart)} / {Formatter('date', item.dateEnd)}
							</TableCell>
							<TableCell width={850}>
								{item.notesStart && <CellText value={item.notesStart} />}
								{item.notesStart && item.notesEnd ? '/' : ''}
								{item.notesEnd && <CellText value={item.notesEnd} />}
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Box>
		</Table>
	)
}
