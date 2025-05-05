import { FC } from 'react'
import { Button, Stack } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'
import { ReactSortable, SortableEvent } from 'react-sortablejs'
import { toast } from 'react-toastify'
import './style.css'

import type { IFetchError } from '@/app/types/error'
import type { ICreateFormStep } from '../../types/create'
import { useAppDispatch } from '@/hooks/redux'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Fallback } from '@/components/Fallback/Fallback'
import { useGetCreateFormStepsQuery, useUpdateSeveralFieldsToCreateFormMutation } from '../../formApiSlice'
import { GroupDialog } from '../Dialogs/GroupDialog'
import { FieldDialog } from '../Dialogs/FieldDialog'
import { Group } from './Group'
import { sortableOptions } from './options'

type Props = {
	section: string
}

export const FieldsList: FC<Props> = ({ section }) => {
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetCreateFormStepsQuery(section, { skip: !section })
	const [updateAll, { isLoading }] = useUpdateSeveralFieldsToCreateFormMutation()

	const methods = useForm<{ data: ICreateFormStep[] }>({ values: { data: data?.data || [] } })
	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = methods
	const { fields, move, update } = useFieldArray({ control, name: 'data' })

	const openDialog = () => {
		const context = { item: { step: data?.data.length }, section: section }
		dispatch(changeDialogIsOpen({ variant: 'CreateFormGroup', isOpen: true, context }))
	}

	const dropHandler = (event: SortableEvent) => {
		console.log('event', event)
		if (event.oldIndex == undefined || event.newIndex == undefined) return
		if (!event.pullMode) move(event.oldIndex, event.newIndex)
	}

	const submitHandler = handleSubmit(async form => {
		const tmp = form.data.map((d, idx) => {
			d.step = idx
			d?.fields?.forEach((item, index) => {
				item.position = index
				item.step = d.step
				item.stepName = d.stepName
			})
			return d.fields || []
		})
		const newData = tmp.flat()
		console.log(newData)

		try {
			await updateAll(newData).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} onSubmit={submitHandler}>
			{isFetching || isLoading ? <Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} /> : null}

			<Stack direction={'row'} justifyContent={'space-between'} mb={1} mx={2}>
				<Button onClick={openDialog} variant='outlined' sx={{ width: 160, textTransform: 'inherit' }}>
					Новая группа
				</Button>

				<Button
					type={'submit'}
					disabled={!Object.keys(dirtyFields).length}
					variant='outlined'
					sx={{ width: 160, textTransform: 'inherit' }}
				>
					Сохранить
				</Button>
			</Stack>

			<FormProvider {...methods}>
				<ReactSortable list={fields} setList={() => {}} onEnd={dropHandler} handle='.drag' {...sortableOptions}>
					{fields.map((item, idx) => (
						<Group key={item.id} index={idx} data={item} section={section} update={update} />
					))}
				</ReactSortable>
			</FormProvider>

			<GroupDialog />
			<FieldDialog />
		</Stack>
	)
}
