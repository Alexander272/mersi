import { CSSProperties, FC } from 'react'

import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import type { ISI } from '../../types/si'
import { useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { Formatter } from '../../utils/formatter'
import { TableRow } from '@/components/Table/TableRow'
import { TableCell } from '@/components/Table/TableCell'
import { CellText } from '@/components/CellText/CellText'

type Props = {
	item: ISI
	sx?: CSSProperties
}

export const Row: FC<Props> = ({ item, sx }) => {
	const section = useAppSelector(getSection)

	const { data } = useGetColumnsQuery(section?.id || '', { skip: !section?.id })

	return (
		<TableRow sx={{ padding: '0 6px', ...sx }}>
			{data?.data.map(c => {
				if (c.children) {
					return c.children.map(c => <Cell key={c.id} item={item} col={c} />)
				}
				return <Cell key={c.id} item={item} col={c} />
			})}
		</TableRow>
	)
}

type CellProps = {
	item: ISI
	col: IColumn
}

const Cell: FC<CellProps> = ({ item, col }) => {
	return (
		<TableCell key={item.id + col.field} width={col.width}>
			<CellText value={Formatter(col.type, item[col.field as keyof ISI])} />
		</TableCell>
	)
}
