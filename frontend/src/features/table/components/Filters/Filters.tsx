import { useRef, useState } from 'react'
import { Badge, Box, Button, Stack, Tooltip, Typography, useTheme } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getFilters, setFilters } from '../../tableSlice'
import { Popover } from '@/components/Popover/Popover'
import { FilterIcon } from '@/components/Icons/FilterIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { CheckIcon } from '@/components/Icons/CheckSimpleIcon'
import { Tabs } from './Tabs'
import { Default } from './Default'
import { Custom } from './Custom'

const defaultValue = {
	field: 'name',
	fieldType: 'text' as const,
	compareType: 'con' as const,
	value: '',
}

export const Filters = () => {
	const [open, setOpen] = useState(false)
	const [filter, setFiler] = useState<'default' | 'custom'>('default')
	const anchor = useRef(null)

	const { palette } = useTheme()

	const filters = useAppSelector(getFilters)
	const dispatch = useAppDispatch()

	const methods = useForm<{ filters: IFilter[] }>()
	const fieldsMethods = useFieldArray({ control: methods.control, name: 'filters' })

	const toggleHandler = () => setOpen(prev => !prev)

	const filterHandler = (value: string) => {
		setFiler(value as 'default')
	}

	const resetHandler = () => {
		fieldsMethods.remove()
		dispatch(setFilters([]))
		toggleHandler()
	}

	const addNewHandler = () => {
		fieldsMethods.append(defaultValue)
	}

	const applyHandler = methods.handleSubmit(form => {
		console.log('form', form)

		const groupedMap = new Map<string, IFilter[]>()
		for (const e of form.filters) {
			let thisList = groupedMap.get(e.field)
			if (thisList === undefined) {
				thisList = []
				groupedMap.set(e.field, thisList)
			}
			thisList.push(e)
		}
		const filters: IFilter[] = []
		groupedMap.forEach(v => filters.push(...v))

		dispatch(setFilters(filters))
		toggleHandler()
	})

	return (
		<>
			<Button
				ref={anchor}
				onClick={toggleHandler}
				variant='outlined'
				color='inherit'
				sx={{ minWidth: 30, paddingX: 1.5 }}
			>
				<Badge color='primary' variant={filters.length < 2 ? 'dot' : 'standard'} badgeContent={filters.length}>
					<FilterIcon fontSize={20} />
				</Badge>
				{/* Фильтр */}
			</Button>

			<Popover
				open={open}
				onClose={toggleHandler}
				anchorEl={anchor.current}
				paperSx={{
					maxWidth: filter == 'default' ? 500 : 700,
					transition: 'all 0.3s ease-in-out',
					'&:before': {
						content: '""',
						display: 'block',
						position: 'absolute',
						top: 0,
						right: filter == 'default' ? '41%' : '29.3%',
						width: 10,
						height: 10,
						bgcolor: 'background.paper',
						transform: 'translate(-50%, -50%) rotate(45deg)',
						zIndex: 0,
					},
				}}
			>
				<Box>
					<Stack direction={'row'} mb={0.5} justifyContent={'space-between'} alignItems={'center'}>
						<Typography fontSize={'1.1rem'}>Фильтр</Typography>

						<Stack direction={'row'} spacing={1} height={34}>
							{filter != 'default' && (
								<Button
									onClick={addNewHandler}
									variant='outlined'
									sx={{ minWidth: 40, padding: '5px 14px' }}
								>
									<PlusIcon fill={palette.primary.main} fontSize={14} />
								</Button>
							)}

							<Tooltip title='Сбросить фильтры' enterDelay={700}>
								<Button onClick={resetHandler} variant='outlined' color='inherit' sx={{ minWidth: 40 }}>
									<TimesIcon fontSize={12} />
								</Button>
							</Tooltip>

							<Tooltip title='Применить фильтры' enterDelay={700}>
								<Button
									onClick={applyHandler}
									variant='contained'
									sx={{ minWidth: 40, padding: '6px 12px' }}
								>
									<CheckIcon fill={palette.common.white} fontSize={20} />
								</Button>
							</Tooltip>
						</Stack>
					</Stack>
					<Tabs value={filter} onChange={filterHandler} />

					<FormProvider {...methods}>
						{filter == 'default' ? <Default /> : <Custom methods={fieldsMethods} />}
					</FormProvider>
				</Box>
			</Popover>
		</>
	)
}
