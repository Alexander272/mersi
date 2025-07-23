import { useRef, useState } from 'react'
import { Badge, Button, ListItemIcon, Menu, MenuItem } from '@mui/material'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { ISort, IFilter } from '../../types/params'
import { PermRules } from '@/constants/permissions'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetSI } from '../../hooks/getSI'
import { useLazyGetSIQuery } from '../../siApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getFilters, getSearch, getSelected, getSort, getStatus, setSelected } from '../../tableSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { CheckListIcon } from '@/components/Icons/CheckListIcon'
import { CheckAllIcon } from '@/components/Icons/CheckAllIcon'
import { DelayIcon } from '@/components/Icons/DelayIcon'
import { CalendarIcon } from '@/components/Icons/CalendarIcon'

export const FastSelect = () => {
	const [open, setOpen] = useState(false)
	const anchor = useRef(null)

	const section = useAppSelector(getSection)
	const selected = useAppSelector(getSelected)
	const status = useAppSelector(getStatus)
	const sort = useAppSelector(getSort)
	const filters = useAppSelector(getFilters)
	const search = useAppSelector(getSearch)
	const all = useCheckPermission(PermRules.Location.Write)
	const dispatch = useAppDispatch()

	const [fetchSi, { isFetching }] = useLazyGetSIQuery()
	const { data } = useGetSI()

	const toggleHandler = () => setOpen(prev => !prev)

	const selectAllHandler = () => {
		if (Object.keys(selected).length) {
			dispatch(setSelected([]))
			toggleHandler()
		} else fetching(filters, sort)
	}

	const selectOverdueHandler = () => {
		const newFilter: IFilter = {
			field: 'nextVerificationDate',
			fieldType: 'date',
			compareType: 'lte',
			value: dayjs().startOf('d').unix().toString(),
		}

		fetching(filters ? [...filters, newFilter] : [newFilter])
	}

	const selectMonthHandler = (value: number) => () => {
		const newFilter: IFilter[] = [
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'gte',
				value: dayjs().add(value, 'M').startOf('month').unix().toString(),
			},
			{
				field: 'nextVerificationDate',
				fieldType: 'date',
				compareType: 'lte',
				value: dayjs().add(value, 'M').endOf('month').unix().toString(),
			},
		]

		fetching(filters ? [...filters, ...newFilter] : [...newFilter])
	}

	const fetching = async (filters?: IFilter[], sort?: ISort) => {
		try {
			const params = {
				all,
				status,
				size: data?.total,
				filters,
				sort,
				search,
				section: section?.id || '',
			}

			const payload = await fetchSi(params).unwrap()
			const newData = payload.data.map(v => v.id)
			dispatch(setSelected(newData))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		} finally {
			toggleHandler()
		}
	}

	return (
		<>
			<Button
				ref={anchor}
				onClick={toggleHandler}
				variant='outlined'
				color='inherit'
				sx={{ minWidth: 30, paddingX: 1.5 }}
			>
				<Badge
					color='primary'
					variant={Object.keys(selected).length < 2 ? 'dot' : 'standard'}
					badgeContent={Object.keys(selected).length}
				>
					<CheckListIcon fontSize={18} />
				</Badge>
				{/* Выбрать */}
			</Button>

			<Menu
				open={open}
				onClose={toggleHandler}
				anchorEl={anchor.current}
				transformOrigin={{ horizontal: 'center', vertical: 'top' }}
				anchorOrigin={{ horizontal: 'center', vertical: 'bottom' }}
				slotProps={{
					paper: {
						elevation: 0,
						sx: {
							overflow: 'visible',
							filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
							mt: 1.5,
							maxWidth: 300,
							width: '100%',
							'&:before': {
								content: '""',
								display: 'block',
								position: 'absolute',
								top: 0,
								right: '47.5%',
								width: 10,
								height: 10,
								bgcolor: 'background.paper',
								transform: 'translate(-50%, -50%) rotate(45deg)',
								zIndex: 0,
							},
						},
					},
				}}
			>
				{isFetching && <BoxFallback />}

				<MenuItem
					onClick={selectAllHandler}
					sx={{ fontWeight: Object.keys(selected).length ? 'bold' : 'normal' }}
				>
					<ListItemIcon>
						<CheckAllIcon fontSize={18} fill={'#474747'} />
					</ListItemIcon>
					{Object.keys(selected).length ? 'Отменить выбор' : 'Выбрать все инструменты'}
				</MenuItem>

				<MenuItem onClick={selectOverdueHandler}>
					<ListItemIcon>
						<DelayIcon fontSize={20} fill={'#474747'} />
					</ListItemIcon>
					Все инструменты у которых
					<br />
					срок поверки прошел
				</MenuItem>

				<MenuItem onClick={selectMonthHandler(0)}>
					<ListItemIcon>
						<CalendarIcon fontSize={18} fill='#474747' />
					</ListItemIcon>
					Все инструменты со сроком <br />
					поверки в текущем месяце
				</MenuItem>

				<MenuItem onClick={selectMonthHandler(1)}>
					<ListItemIcon>
						<CalendarIcon fontSize={18} fill='#474747' />
					</ListItemIcon>
					Все инструменты со сроком
					<br /> поверки в следующем месяце
				</MenuItem>
			</Menu>
		</>
	)
}
