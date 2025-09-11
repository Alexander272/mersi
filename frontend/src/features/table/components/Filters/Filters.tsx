import { useEffect, useRef, useState } from 'react'
import { Badge, Box, Button, Popover, Stack, Tooltip, Typography, useTheme } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IFilter } from '../../types/params'
import { localKeys } from '@/constants/localKeys'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useSaveFiltersMutation } from '../../filtersApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getFilters, setFilters } from '../../tableSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
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
	const [filter, setFilter] = useState<'default' | 'custom'>('default')
	const anchor = useRef(null)

	const { palette } = useTheme()

	const filters = useAppSelector(getFilters)
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const [save, { isLoading }] = useSaveFiltersMutation()

	const methods = useForm<{ filters: IFilter[] }>()
	const fieldsMethods = useFieldArray({ control: methods.control, name: 'filters' })

	useEffect(() => {
		if (filters) {
			methods.reset({ filters })
		}
	}, [filters, methods])

	const toggleHandler = () => setOpen(prev => !prev)

	const filterHandler = (value: string) => {
		setFilter(value as 'default')
	}

	const resetHandler = async () => {
		fieldsMethods.remove()
		try {
			await save({ filters: [], section: section?.id || '' }).unwrap()
			localStorage.removeItem(localKeys.activeFilters)
			dispatch(setFilters([]))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		} finally {
			toggleHandler()
		}
	}

	const addNewHandler = () => {
		fieldsMethods.append(defaultValue)
	}

	const applyHandler = methods.handleSubmit(async form => {
		console.log('form', form)

		const groupedMap = new Map<string, IFilter[]>()
		if (form.filters?.length) {
			for (const e of form.filters) {
				let thisList = groupedMap.get(e.field)
				if (thisList === undefined) {
					thisList = []
					groupedMap.set(e.field, thisList)
				}
				thisList.push(e)
			}
		}
		const filters: IFilter[] = []
		groupedMap.forEach(v => filters.push(...v))

		try {
			await save({ filters: filters, section: section?.id || '' }).unwrap()
			dispatch(setFilters(filters))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		} finally {
			toggleHandler()
		}
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
				anchorOrigin={{
					vertical: 'bottom',
					horizontal: 'center',
				}}
				transformOrigin={{
					vertical: 'top',
					horizontal: 'center',
				}}
				slotProps={{
					paper: {
						elevation: 0,
						sx: {
							overflow: 'visible',
							filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
							mt: 1,
							paddingX: 2,
							pt: 1.5,
							paddingBottom: 2,
							maxWidth: filter == 'default' ? 500 : 750,
							width: '100%',
							'&:before': {
								content: '""',
								display: 'block',
								position: 'absolute',
								top: 0,
								right: filter == 'default' ? '41%' : '27.5%',
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
				<Box>
					{isLoading && <BoxFallback />}

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
