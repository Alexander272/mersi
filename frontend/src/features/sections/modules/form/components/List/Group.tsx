import { FC } from 'react'
import { IconButton, Stack, Typography } from '@mui/material'
import { useFieldArray, UseFieldArrayUpdate, useFormContext } from 'react-hook-form'
import { ReactSortable, SortableEvent } from 'react-sortablejs'
import { toast } from 'react-toastify'

import type { ICreateFormStep } from '../../types/create'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { DragIcon } from '@/components/Icons/DragIcon'
import { EditIcon } from '@/components/Icons/EditIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { sortableOptions } from './options'
import { Row } from './Row'

type Props = {
	index: number
	data: ICreateFormStep
	section: string
	update: UseFieldArrayUpdate<{ data: ICreateFormStep[] }, 'data'>
}

export const Group: FC<Props> = ({ data, section, index, update }) => {
	const dispatch = useAppDispatch()

	const { control, getValues } = useFormContext<{ data: ICreateFormStep[] }>()
	const { fields, move, remove } = useFieldArray({ control, name: `data.${index}.fields` })

	const openEditDialog = () => {
		const context = { item: data, section }
		dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: true, context }))
	}

	const openDialog = () => {
		const context = {
			item: {
				step: data?.step || 0,
				stepName: data.stepName,
				sectionId: section || '',
				position: data.fields.length,
			},
		}
		dispatch(changeDialogIsOpen({ variant: 'CreateFormField', isOpen: true, context }))
	}

	const dropHandler = (event: SortableEvent) => {
		console.log('event', event)
		if (!event || event.oldIndex == undefined || event.newIndex == undefined) return
		if (!event.pullMode) {
			move(event.oldIndex, event.newIndex)
			return
		}

		if (!data.fields) return
		const values = getValues()
		const idx = values.data.findIndex(d => d.stepName == event.to.id)
		if (idx == -1) {
			toast.error('Не удалось найти группу для переноса')
			return
		}
		const item = { ...values.data[idx] }
		const old = values.data.find(d => d.stepName == event.from.id)
		if (!old) return
		old.step = item.step
		old.stepName = item.stepName
		item.fields.splice(event.newIndex, 0, old.fields[event.oldIndex])

		update(idx, item)
		remove(event.oldIndex)
	}

	return (
		<Stack border={'2px solid #f5f5f5'} borderRadius={3} mb={2}>
			<Stack direction={'row'} alignItems={'center'} position={'relative'}>
				<IconButton sx={{ cursor: 'grab', mr: 1 }} className='drag'>
					<DragIcon fill={'#a8a8a8'} fontSize={24} />
				</IconButton>

				<Typography sx={{ ml: 'auto' }}>{data.stepName}</Typography>
				<IconButton onClick={openEditDialog} size='large' sx={{ ml: 0.5, mr: 'auto' }}>
					<EditIcon fontSize={12} fill={'#6e6e6e'} />
				</IconButton>

				<IconButton onClick={openDialog} size='large' sx={{ mr: 2 }}>
					<PlusIcon fontSize={10} fill={'#575757'} />
				</IconButton>
			</Stack>

			<Stack sx={{ pl: 1.5, pt: 1, background: '#f5f5f5', borderRadius: 3 }}>
				<ReactSortable
					id={data.stepName}
					list={data.fields}
					setList={() => {}}
					onEnd={dropHandler}
					handle='.drag'
					{...sortableOptions}
				>
					{fields.map((item, idx) => (
						<Row key={item.id} id={data.fields[idx].id} data={item} />
					))}
				</ReactSortable>
			</Stack>
		</Stack>
	)
}
