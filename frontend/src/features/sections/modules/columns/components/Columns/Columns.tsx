import { FC } from 'react'
import { Button, Stack } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'
import { ReactSortable } from 'react-sortablejs'
import { SortableEvent } from 'sortablejs'
import './style.css'

import type { IFetchError } from '@/app/types/error'
import type { IColumn } from '../../types/columns'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { useGetColumnsQuery, useUpdateColumnPositionsMutation } from '../../columnsApiSlice'
import { ColumnDialog } from '../Dialog/Dialog'
import { Column } from './Column'
import { sortableOptions } from './options'

type Props = {
	section: string
}

export const Columns: FC<Props> = ({ section }) => {
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetColumnsQuery(section, { skip: !section })
	const [updateColumns, { isLoading }] = useUpdateColumnPositionsMutation()

	const methods = useForm<{ data: IColumn[] }>({ values: { data: data?.data || [] } })
	const {
		control,
		handleSubmit,
		// formState: { dirtyFields },
	} = methods
	const { fields, move, insert, remove, update } = useFieldArray({ control, name: 'data' })

	const openDialog = () => {
		const context = {
			item: { sectionId: section, position: (data?.data[data?.data.length - 1]?.position || 0) + 1 },
		}
		dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: true, context }))
	}

	const updateHandler = handleSubmit(async form => {
		let index = 1
		const tmp = form.data.map(d => {
			d.position = index
			index++
			d?.children?.forEach(item => {
				item.position = index
				index++
				item.parentId = d.id
			})
			return [d, ...(d.children || [])]
		})
		const newData = tmp.flat()
		console.log(newData)

		try {
			await updateColumns(newData).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const dropHandler = (event: SortableEvent) => {
		console.log('event', event)
		if (event.oldIndex == undefined || event.newIndex == undefined) return
		// update(event.newIndex, { ...fields[event.newIndex], position: event.newIndex + 1 })
		// update(event.oldIndex, { ...fields[event.oldIndex], position: event.oldIndex + 1 })
		if (!event.pullMode) {
			move(event.oldIndex, event.newIndex)
			return
		}

		let idx = fields.findIndex(item => item.id == event.to.id)
		if (idx == -1 || !data) return
		const item = { ...fields[idx] }
		item.children?.splice(event.newIndex, 0, fields[event.oldIndex])
		remove(event.oldIndex)
		if (idx > event.oldIndex) idx--
		update(idx, item)
	}

	return (
		<Stack component={'form'} onSubmit={updateHandler}>
			{isFetching || isLoading ? <Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} /> : null}

			<Stack direction={'row'} justifyContent={'space-between'} mb={1} mx={2}>
				<Button
					// onClick={updateHandler}
					type='submit'
					variant='outlined'
					// disabled={!Object.keys(dirtyFields).length}
					sx={{ textTransform: 'inherit' }}
				>
					Обновить индексы
				</Button>

				<Button onClick={openDialog} variant='outlined' sx={{ width: 160, textTransform: 'inherit' }}>
					Новая
				</Button>
			</Stack>

			<FormProvider {...methods}>
				<ReactSortable list={fields} setList={() => {}} onEnd={dropHandler} handle='.drag' {...sortableOptions}>
					{fields.map((item, idx) => (
						<Column key={item.id} index={idx} id={data?.data[idx]?.id || ''} data={item} insert={insert} />
					))}
				</ReactSortable>
			</FormProvider>

			<ColumnDialog />
		</Stack>
	)
}
