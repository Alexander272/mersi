import { useEffect, useRef, useState } from 'react'
import { Button, IconButton, Stack, Tooltip, Typography, useTheme } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'

import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import type { IChangedColumns } from '../../types/table'
import { localKeys } from '@/constants/localKeys'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { apiSlice } from '@/app/apiSlice'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getChangedColumns, setChangedColumns } from '../../tableSlice'
import { Popover } from '@/components/Popover/Popover'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { SettingIcon } from '@/components/Icons/SettingIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { CheckIcon } from '@/components/Icons/CheckSimpleIcon'
import { Item } from './Item'
import { Group } from './Group'

export const Setting = () => {
	const { palette } = useTheme()
	const section = useAppSelector(getSection)
	const changedColumns = useAppSelector(getChangedColumns)
	const dispatch = useAppDispatch()

	const [open, setOpen] = useState(false)
	const anchor = useRef(null)

	const { data, isFetching } = useGetColumnsQuery(
		{ section: section?.id || '', original: true },
		{ skip: !section?.id }
	)

	// const changedData: IColumn[] | undefined = data?.data.map(c => {
	// 	if (!c.children?.length) return { ...c, ...changedColumns?.[c.id] }
	// 	const newChildren = c.children.map(c => ({ ...c, ...changedColumns?.[c.id] }))
	// 	return { ...c, ...changedColumns?.[c.id], children: newChildren }
	// })
	const methods = useForm<{ data: IColumn[] }>({
		values: { data: data?.data || [] },
	})
	const { control, handleSubmit } = methods
	const { fields } = useFieldArray({ control, name: 'data' })

	useEffect(() => {
		const changedData: IColumn[] | undefined = data?.data.map(c => {
			if (!c.children?.length) return { ...c, ...changedColumns?.[c.id] }
			const newChildren = c.children.map(c => ({ ...c, ...changedColumns?.[c.id] }))
			return { ...c, ...changedColumns?.[c.id], children: newChildren }
		})

		methods.setValue('data', changedData || [])
	}, [changedColumns, data, methods])

	const toggleHandler = () => setOpen(prev => !prev)

	const resetHandler = () => {
		if (!section) return
		dispatch(setChangedColumns(undefined))
		localStorage.removeItem(localKeys.changedColumns(section.id))
		toggleHandler()
		dispatch(apiSlice.util.invalidateTags([{ type: 'Columns', id: 'List' }]))
	}

	const applyHandler = handleSubmit(form => {
		// console.log('apply', form)
		// console.log('dirty', methods.formState.dirtyFields)
		if (!section) return

		const columns: IChangedColumns = {}
		methods.formState.dirtyFields.data?.forEach((d, i) => {
			if (!d.children)
				columns[form.data[i].id] = {
					id: form.data[i].id,
					hidden: form.data[i].hidden || false,
					width: +form.data[i].width,
				}
			else {
				columns[form.data[i].id] = {
					id: form.data[i].id,
					hidden: form.data[i].hidden || false,
					width: +form.data[i].width,
				}
				d.children.forEach((_c, j) => {
					const data = form.data[i].children?.[j]
					if (data) columns[data.id] = { id: data.id, hidden: data.hidden || false, width: +data.width }
				})
			}
		})
		dispatch(setChangedColumns(columns))
		localStorage.setItem(localKeys.changedColumns(section.id), JSON.stringify(columns))
		toggleHandler()
		dispatch(apiSlice.util.invalidateTags([{ type: 'Columns', id: 'List' }]))
	})

	return (
		<>
			<IconButton ref={anchor} onClick={toggleHandler}>
				<SettingIcon fontSize={20} />
			</IconButton>

			<Popover
				open={open}
				onClose={toggleHandler}
				anchorEl={anchor.current}
				paperSx={{ padding: 0, maxWidth: 500 }}
			>
				<Stack>
					{isFetching && <BoxFallback />}

					<Stack
						direction={'row'}
						mx={2}
						mt={1}
						mb={2.5}
						justifyContent={'space-between'}
						alignItems={'center'}
					>
						<Typography fontSize={'1.1rem'}>Настройка колонок</Typography>

						<Stack direction={'row'} spacing={1} height={34}>
							<Tooltip title='Сбросить настройки'>
								<Button onClick={resetHandler} variant='outlined' color='inherit' sx={{ minWidth: 40 }}>
									<RefreshIcon fontSize={18} />
								</Button>
							</Tooltip>

							<Button
								onClick={applyHandler}
								variant='contained'
								sx={{ minWidth: 40, padding: '6px 15px' }}
							>
								<CheckIcon fill={palette.common.white} fontSize={20} />
							</Button>
						</Stack>
					</Stack>

					<Stack maxHeight={450} overflow={'auto'} ml={1.5} pr={1}>
						<FormProvider {...methods}>
							{fields.map((f, i) => {
								if (f.type == 'parent') return <Group key={f.id} index={i} data={f} />
								return <Item key={f.id} index={i} label={f.name} />
							})}
						</FormProvider>
					</Stack>
				</Stack>
			</Popover>
		</>
	)
}

export default Setting
