import { FC } from 'react'
import { IconButton, Stack, Typography } from '@mui/material'
import { useFieldArray, UseFieldArrayInsert, useFormContext } from 'react-hook-form'
import { ReactSortable, SortableEvent } from 'react-sortablejs'

import type { IColumn } from '../../types/columns'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { DragIcon } from '@/components/Icons/DragIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { EditIcon } from '@/components/Icons/EditIcon'
import { sortableOptions } from './options'

type Props = {
	id: string
	index: number
	data: IColumn
	isChild?: boolean

	insert: UseFieldArrayInsert<{ data: IColumn[] }>
}

export const Column: FC<Props> = ({ index, id, data, isChild, insert }) => {
	if (!data.children) return <Row index={index} id={id} data={data} isChild={isChild} />
	return <Group index={index} id={id} data={data} isChild={isChild} insert={insert} />
}

const Row: FC<Omit<Props, 'insert'>> = ({ id, data, isChild }) => {
	const dispatch = useAppDispatch()

	const openEditDialog = () => {
		dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: true, context: { item: { ...data, id } } }))
	}
	const openDialog = () => {
		const context = { item: { parentId: id, sectionId: data.sectionId, position: data.position } }
		dispatch(changeDialogIsOpen({ variant: 'Columns', isOpen: true, context }))
	}

	return (
		<Stack
			direction={'row'}
			alignItems={'center'}
			mb={1}
			position={'relative'}
			sx={{
				':after': {
					content: '""',
					position: 'absolute',
					left: 40,
					bottom: 0,
					height: '1px',
					width: 'calc(100% - 60px)',
					background: '#a0a6b7a3',
				},
			}}
		>
			<IconButton sx={{ cursor: 'grab', mr: 1 }} className='drag'>
				<DragIcon fill={'#a8a8a8'} fontSize={24} />
			</IconButton>

			<Typography>{data.name}</Typography>
			<IconButton onClick={openEditDialog} size='large' sx={{ ml: 0.5 }}>
				<EditIcon fontSize={12} fill={'#6e6e6e'} />
			</IconButton>

			{!isChild && (
				<IconButton onClick={openDialog} size='large' sx={{ ml: 'auto', mr: 3 }}>
					<PlusIcon fontSize={10} fill={'#575757'} />
				</IconButton>
			)}
		</Stack>
	)
}

const Group: FC<Props> = ({ index, id, data, isChild, insert }) => {
	const { control } = useFormContext<{ data: IColumn[] }>()
	const { fields, move, remove } = useFieldArray({ control, name: `data[${index}].children` as 'data' })

	const dropHandler = (event: SortableEvent) => {
		console.log('event', event)
		if (!event || event.oldIndex == undefined || event.newIndex == undefined) return
		if (!event.pullMode) {
			move(event.oldIndex, event.newIndex)
			return
		}

		if (!data.children) return
		remove(event.oldIndex)
		insert(event.newIndex, data.children[event.oldIndex])
	}

	return (
		<Stack pb={1}>
			<Row index={index} id={id} data={data} isChild={isChild} />
			{data.children && (
				<Stack sx={{ mx: 2, pl: 2, pt: 1, background: '#f5f5f5', borderRadius: 3 }}>
					<ReactSortable
						id={data.id}
						list={data.children}
						setList={() => {}}
						onEnd={dropHandler}
						handle='.drag'
						{...sortableOptions}
					>
						{fields.map((item, idx) => (
							<Column key={item.id} index={idx} id={id} data={item} isChild insert={insert} />
						))}
					</ReactSortable>
				</Stack>
			)}
		</Stack>
	)
}
