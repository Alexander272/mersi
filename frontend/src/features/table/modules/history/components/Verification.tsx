import { FC } from 'react'
import { Box, Chip, Stack } from '@mui/material'
import dayjs from 'dayjs'

import type { IVerDocs } from '../../verification/types/verificationDocs'
import { useAppSelector } from '@/hooks/redux'
import { useLazyDownloadFileQuery } from '@/features/files/fileApiSlice'
import { useGetVerificationFieldsQuery, useGetVerificationsQuery } from '../../verification/verificationApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { Table } from '@/components/Table/Table'
import { TableHead } from '@/components/Table/TableHead'
import { TableBody } from '@/components/Table/TableBody'
import { TableRow } from '@/components/Table/TableRow'
import { TableCell } from '@/components/Table/TableCell'
import { CellText } from '@/components/CellText/CellText'
import { DocIcon } from '@/components/Icons/DocIcon'
import { PdfIcon } from '@/components/Icons/PdfIcon'
import { ImageIcon } from '@/components/Icons/ImageIcon'
import { SheetIcon } from '@/components/Icons/SheetIcon'

const FileTypes = {
	doc: <DocIcon ml={0.8} />,
	pdf: <PdfIcon ml={0.8} />,
	image: <ImageIcon ml={0.8} />,
	sheet: <SheetIcon ml={0.8} />,
}

const widths = new Map([
	[
		'46ba9e17-65c7-474b-8c47-7975ab4319d5',
		{ verificationDate: 180, nextVerificationDate: 180, status: 160, notes: 315, docs: 315 },
	],
	[
		'5f86ac32-9477-4f48-ae27-e9ef94d848f8',
		{ verificationDate: 240, nextVerificationDate: 240, notes: 335, docs: 335 },
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
	const [download] = useLazyDownloadFileQuery()

	const downloadHandler = (doc: IVerDocs) => async () => {
		if (doc.docId == '') return
		await download({ path: doc.path, label: doc.doc })
	}

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
				</TableRow>
			</TableHead>
			<Box maxHeight={350} overflow={'auto'} position={'relative'} minHeight={150}>
				<TableBody>
					{data?.data.map(item => (
						<TableRow key={item.id} sx={{ minHeight: 38, cursor: 'default' }}>
							{fields?.data.map(f => {
								const sizes = widths.get(section?.id || '')

								if (f.field == 'docs') {
									return (
										<TableCell key={f.id} width={sizes?.[f.field as 'notes']}>
											<Stack direction={'row'} justifyContent={'center'} alignItems={'center'}>
												{item.docs?.map(d => (
													<Chip
														key={d.id}
														icon={FileTypes[d.type as 'doc']}
														label={d.doc}
														onClick={downloadHandler(d)}
														variant='outlined'
														clickable={d.docId != ''}
													/>
												))}
											</Stack>
										</TableCell>
									)
								}

								return (
									<TableCell key={f.id} width={sizes?.[f.field as 'notes']}>
										<CellText value={Formatter(f.field, item[f.field as 'notes'])} />
									</TableCell>
								)
							})}
						</TableRow>
					))}
				</TableBody>
			</Box>
		</Table>
	)
}
