import { FC } from 'react'
import { Box } from '@mui/material'
import dayjs from 'dayjs'

import { useAppSelector } from '@/hooks/redux'
import { useGetVerificationFieldsQuery, useGetVerificationsQuery } from '../../verification/verificationApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Table } from '@/components/Table/Table'
import { TableHead } from '@/components/Table/TableHead'
import { TableBody } from '@/components/Table/TableBody'
import { TableRow } from '@/components/Table/TableRow'
import { TableCell } from '@/components/Table/TableCell'
import { CellText } from '@/components/CellText/CellText'

const widths = new Map([
	[
		'46ba9e17-65c7-474b-8c47-7975ab4319d5',
		{ verificationDate: 180, nextVerificationDate: 180, status: 160, notes: 300, 'docs.0.doc': 300 },
	],
])

const Statuses = {
	work: 'Пригоден',
	repair: 'Нужен ремонт',
	decommissioning: 'Не пригоден',
}

const Formatter = (field: string, value: string) => {
	if (field == 'verificationDate' || field == 'nextVerificationDate') return dayjs(+value * 1000).format('DD.MM.YYYY')
	if (field == 'status') return Statuses[value as 'work' | 'repair' | 'decommissioning']
	return value
}

type Props = {
	instrumentId: string
}

export const Verification: FC<Props> = ({ instrumentId }) => {
	const section = useAppSelector(getSection)

	const { data: fields, isFetching: isFetchingFields } = useGetVerificationFieldsQuery(
		{ section: section?.id || '', group: 'history' },
		{ skip: !section?.id }
	)
	const { data, isFetching } = useGetVerificationsQuery(instrumentId, { skip: !instrumentId })

	if (isFetching || isFetchingFields) return <BoxFallback />
	return (
		<Table>
			<TableHead>
				<TableRow height={80}>
					{fields?.data.map(item => {
						const sizes = widths.get(section?.id || '')
						return (
							<TableCell key={item.id} width={sizes?.[item.field as 'notes']}>
								{item.label}
							</TableCell>
						)
					})}
					{/* <TableCell>Дата поверки (калибровки)</TableCell> */}
					{/* //TODO дописать */}
				</TableRow>
			</TableHead>
			<Box maxHeight={350} overflow={'auto'} position={'relative'} minHeight={150}>
				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id} height={38}>
							{fields?.data.map(f => {
								const sizes = widths.get(section?.id || '')
								return (
									<TableCell key={f.id} width={sizes?.[f.field as 'notes']}>
										<CellText value={Formatter(f.field, item[f.field as 'notes'])} />
									</TableCell>
								)
							})}
							{/* <TableCell>{Formatter('date', item.verificationDate)}</TableCell> */}
						</TableRow>
					))}
				</TableBody>
			</Box>
		</Table>
	)
}
