import { JSX } from 'react'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { IColumn } from '@/features/sections/modules/columns/types/columns'
import { TableCell } from '@/components/Table/TableCell'
import { TableHead } from '@/components/Table/TableHead'
import { TableRow } from '@/components/Table/TableRow'
import { TableGroup } from '@/components/Table/TableGroup'
import { CellText } from '@/components/CellText/CellText'
import { Badge } from '@/components/Badge/Badge'
import { Fallback } from '@/components/Fallback/Fallback'
import { SortUpIcon } from '@/components/Icons/SortUpIcon'
import { ColWidth, RowHeight } from '../../constants/defaultValues'
import { getHidden, getSort, setSort } from '../../tableSlice'
import { useCalcWidth } from '../../utils/calcWidth'

export const Head = () => {
	const section = useAppSelector(getSection)
	const sort = useAppSelector(getSort)
	const hidden = useAppSelector(getHidden)

	const { data, isFetching } = useGetColumnsQuery(section?.id || '', { skip: !section?.id })

	const { width, hasFewRows } = useCalcWidth(data?.data || [])
	const height = (hasFewRows ? 2 : 1) * RowHeight

	const dispatch = useAppDispatch()

	const setSortHandler = (field: string) => () => {
		dispatch(setSort(field))
	}

	const getCell = (c: IColumn) => {
		return (
			<TableCell
				key={c.field}
				width={c.width || ColWidth}
				isActive
				onClick={c.allowSort ? setSortHandler(c.field) : undefined}
			>
				<CellText value={c.name} />
				{c.allowSort ? (
					<Badge
						color='primary'
						badgeContent={Object.keys(sort).findIndex(k => k == c.field) + 1}
						invisible={Object.keys(sort).length < 2}
					>
						<SortUpIcon
							fontSize={16}
							fill={sort[c.field] ? 'black' : '#adadad'}
							transform={!sort[c.field] || sort[c.field] == 'ASC' ? '' : 'rotateX(180deg)'}
							transition={'.2s all ease-in-out'}
						/>
					</Badge>
				) : null}
			</TableCell>
		)
	}
	const renderHeader = () => {
		const header: JSX.Element[] = []

		data?.data.forEach(c => {
			if (c.children) {
				let width = 0
				const subhead: JSX.Element[] = []

				c.children.forEach(c => {
					if (!hidden[c.field]) {
						width += c.width || ColWidth

						subhead.push(getCell(c))
					}
				})

				if (subhead.length > 0) {
					header.push(
						<TableGroup key={c.field}>
							<TableRow>
								<TableCell width={width} key={c.field}>
									<CellText value={c.name} />
								</TableCell>
							</TableRow>
							<TableRow>{subhead}</TableRow>
						</TableGroup>
					)
				}
			} else if (!hidden[c.field]) {
				header.push(getCell(c))
			}
		})

		return header
	}

	if (isFetching) return <Fallback />
	return (
		<TableHead>
			<TableRow width={width} height={height} sx={{ padding: '0 6px' }}>
				{renderHeader()}
			</TableRow>
		</TableHead>
	)
}
