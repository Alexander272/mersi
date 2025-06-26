import { CSSProperties, FC, MouseEvent } from 'react'
import { useTheme } from '@mui/material'

import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import type { ISI } from '../../types/si'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getContextMenu, getHidden, getSelected, setContextMenu, setSelected } from '../../tableSlice'
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
	const selected = useAppSelector(getSelected)
	const hidden = useAppSelector(getHidden)
	const contextMenu = useAppSelector(getContextMenu)
	const dispatch = useAppDispatch()

	const { palette } = useTheme()

	const { data } = useGetColumnsQuery(section?.id || '', { skip: !section?.id })

	const selectHandler = () => {
		dispatch(setSelected(item.id))
	}

	const contextHandler = (event: MouseEvent<HTMLDivElement>) => {
		event.preventDefault()
		const menu = {
			active: item.id,
			coords: { mouseX: event.clientX + 2, mouseY: event.clientY - 6 },
		}
		dispatch(setContextMenu(menu))
	}

	let background = ''
	if (selected[item.id]) background = palette.rowActive.light
	if (contextMenu?.active == item.id) background = palette.rowActive.main

	return (
		<TableRow onClick={selectHandler} onContext={contextHandler} hover sx={{ padding: '0 6px', ...sx, background }}>
			{data?.data.map(c => {
				if (hidden[c.field]) return null
				if (c.children) {
					if (hidden[c.field]) return null
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
