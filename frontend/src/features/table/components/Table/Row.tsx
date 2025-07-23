import { CSSProperties, FC, MouseEvent } from 'react'
import { useTheme } from '@mui/material'
import dayjs from 'dayjs'

import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import type { ISI } from '../../types/si'
import { Formatter } from '../../utils/formatter'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getRealm } from '@/features/realms/realmSlice'
import { getContextMenu, getHidden, getSelected, getStatus, setContextMenu, setSelected } from '../../tableSlice'
import { TableRow } from '@/components/Table/TableRow'
import { TableCell } from '@/components/Table/TableCell'
import { CellText } from '@/components/CellText/CellText'

const RowColors = {
	overdue: '#ff3f3f',
	deadlineRed: '#ff9393',
	deadlineOr: '#fdae1f',
	moved: '#eee',
	preservation: '#f7f7f7',
	save: '#f7f7f7',
	reverse: '#fff4cb',
}

type Props = {
	item: ISI
	sx?: CSSProperties
}

export const Row: FC<Props> = ({ item, sx }) => {
	const section = useAppSelector(getSection)
	const realm = useAppSelector(getRealm)
	const selected = useAppSelector(getSelected)
	const status = useAppSelector(getStatus)
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

	const getStyles = () => {
		const styles = { background: '' }
		if (status != 'work') return styles

		if (selected[item.id]) styles.background = palette.rowActive.light
		if (contextMenu?.active == item.id) styles.background = palette.rowActive.main

		if (item.status == 'moved') styles.background = RowColors.moved

		if (realm?.id == 'c0d7fd4e-da02-4265-be05-9023bac5ff0d') {
			const deadline = dayjs()
				.add(1, 'month')
				.isAfter(dayjs(item.nextVerificationDate * 1000))
			if (deadline) styles.background = RowColors.deadlineOr
		} else {
			const deadline = dayjs()
				.add(15, 'd')
				.isAfter(dayjs(item.nextVerificationDate * 1000))
			if (deadline) styles.background = RowColors.deadlineRed
		}
		const overdue = dayjs().isAfter(dayjs(item.nextVerificationDate * 1000))
		if (overdue) styles.background = RowColors.overdue

		if (item.status == 'reserve') styles.background = RowColors.reverse

		if (item.preservationDate && !item.dePreservationDate) styles.background = RowColors.preservation
		if (item.transferDate && !item.returnDate) styles.background = RowColors.save

		return styles
	}

	return (
		<TableRow
			onClick={selectHandler}
			onContext={contextHandler}
			hover
			sx={{ padding: '0 6px', ...sx, ...getStyles() }}
		>
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
