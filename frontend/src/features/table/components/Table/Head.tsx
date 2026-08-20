import { JSX, useEffect, useRef } from 'react'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import type { ISort } from '../../types/params'
import { ColWidth, RowHeight } from '../../constants/defaultValues'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useDebounceFunc } from '@/hooks/useDebounceFunc'
import { useCalcWidth } from '../../utils/calcWidth'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { useGetFiltersQuery, useGetSortQuery, useSaveSortMutation } from '../../filtersApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getSort, setChangedColumns, setDefaultSorting, setFilters, setSort } from '../../tableSlice'
import { TableCell } from '@/components/Table/TableCell'
import { TableHead } from '@/components/Table/TableHead'
import { TableRow } from '@/components/Table/TableRow'
import { TableGroup } from '@/components/Table/TableGroup'
import { CellText } from '@/components/CellText/CellText'
import { Badge } from '@/components/Badge/Badge'
import { Fallback } from '@/components/Fallback/Fallback'
import { SortUpIcon } from '@/components/Icons/SortUpIcon'
import { localKeys } from '@/constants/localKeys'

export const Head = () => {
	const saving = useRef(false)
	const section = useAppSelector(getSection)
	const sort = useAppSelector(getSort)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetColumnsQuery({ section: section?.id || '' }, { skip: !section?.id })
	const { data: filters } = useGetFiltersQuery(section?.id || '', { skip: !section?.id })
	const { data: sorting } = useGetSortQuery(section?.id || '', { skip: !section?.id })
	const [save] = useSaveSortMutation()

	useEffect(() => {
		if (filters) dispatch(setFilters(filters.data))
	}, [dispatch, filters])
	useEffect(() => {
		if (sorting)
			dispatch(setDefaultSorting(sorting.data.reduce((acc, v) => ({ ...acc, [v.name]: v.orderType }), {})))
	}, [dispatch, sorting])
	useEffect(() => {
		if (!section) return
		const columns = localStorage.getItem(localKeys.changedColumns(section.id))
		dispatch(setChangedColumns(columns ? JSON.parse(columns) : undefined))
	}, [dispatch, section])

	const { width, hasFewRows } = useCalcWidth(data?.data || [])
	const height = (hasFewRows ? 2 : 1) * RowHeight

	const saveSort = useDebounceFunc(async sort => {
		const data = Object.keys(sort as ISort).map((k, i) => ({ name: k, orderType: (sort as ISort)[k], count: i }))
		try {
			await save({ sort: data, section: section?.id || '' }).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}, 1000)
	useEffect(() => {
		if (!saving.current) return
		saveSort(sort)
		saving.current = false
	}, [saveSort, sort])

	const setSortHandler = (field: string) => () => {
		dispatch(setSort(field))
		saving.current = true
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
			if (c.children && !c?.hidden) {
				let width = 0
				const subhead: JSX.Element[] = []

				c.children.forEach(c => {
					if (!c?.hidden) {
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
			} else if (!c?.hidden) {
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
