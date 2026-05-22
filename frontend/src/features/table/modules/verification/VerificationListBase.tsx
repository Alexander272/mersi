import { FC } from 'react'
import { Button, Stack } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'
import { ReactSortable, SortableEvent } from 'react-sortablejs'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IVerificationFieldDTO, IVerificationFieldForm } from './modules/types/verificationFields'
import { useAppDispatch } from '@/hooks/redux'
import {
	useCreateVerificationFieldsMutation,
	useDeleteVerificationFieldsMutation,
	useGetVerificationFieldsQuery,
	useUpdateVerificationFieldsMutation,
} from './modules/fieldsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

const sortableOptions = {
	animation: 150,
	fallbackOnBody: true,
	swapThreshold: 0.65,
	group: 'columns',
}

type Props = {
	section: string
	group: 'form' | 'history'
	dialogVariant: 'EditVerificationFields' | 'EditVerificationHistory'
	toastMessage: string
	DialogComponent: FC<{ submit: (data: IVerificationFieldForm) => void }>
	ItemComponent: FC<{ data: IVerificationFieldForm }>
}

export const VerificationListBase: FC<Props> = ({
	section,
	group,
	dialogVariant,
	toastMessage,
	DialogComponent,
	ItemComponent,
}) => {
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetVerificationFieldsQuery(
		{ section: section || '', group },
		{ skip: !section || section == 'new' },
	)
	const [create, { isLoading: creating }] = useCreateVerificationFieldsMutation()
	const [updateAll, { isLoading: updating }] = useUpdateVerificationFieldsMutation()
	const [removeAll, { isLoading: removing }] = useDeleteVerificationFieldsMutation()

	const methods = useForm<{ data: IVerificationFieldForm[] }>({
		values: { data: data?.data.map(d => ({ ...d, status: 'none', group: '' })) || [] },
	})
	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = methods
	const { fields, move, update, insert, remove } = useFieldArray({ control, name: 'data', keyName: '_id' })

	const dropHandler = (event: SortableEvent) => {
		if (event.oldIndex == undefined || event.newIndex == undefined || event.oldIndex == event.newIndex) return

		move(event.oldIndex, event.newIndex)
	}

	const openDialog = () => {
		dispatch(changeDialogIsOpen({ variant: dialogVariant, isOpen: true }))
	}

	const submitHandler = (data: IVerificationFieldForm) => {
		if (data.status == 'updated' || data.status == 'deleted') update(data.position - 1, data)
		if (data.status == 'new') {
			if (data.position == 1) insert(0, { ...data, position: 1 })
			else update(data.position - 1, data)
		}
		if (data.status == undefined) remove(data.position - 1)
	}

	const updateHandler = handleSubmit(async form => {
		const updated: IVerificationFieldDTO[] = []
		const created: IVerificationFieldDTO[] = []
		const deleted: string[] = []

		form.data.forEach((item, idx) => {
			if (item.position != idx + 1) {
				item.status = 'moved'
				item.position = idx + 1
			}
			item.sectionId = section
			item.group = group

			if (group == 'history') item.width = +item.width

			if (item.status == 'updated' || item.status == 'moved') updated.push(item)
			if (item.status == 'new') created.push(item)
			if (item.status == 'deleted') deleted.push(item.id)
		})
		if (!updated.length && !created.length && !deleted.length) return

		try {
			if (updated.length) await updateAll(updated).unwrap()
			if (created.length) await create(created).unwrap()
			if (deleted.length) await removeAll(deleted).unwrap()
			methods.reset()
			toast.success(toastMessage)
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<>
			<Stack component={'form'} onSubmit={updateHandler}>
				{isFetching || creating || updating || removing ? <BoxFallback /> : null}

				<Stack direction={'row'} justifyContent={'space-between'} mb={2.5} mx={2}>
					<Button
						onClick={openDialog}
						variant='outlined'
						disabled={section == 'new'}
						sx={{ width: 160, textTransform: 'inherit' }}
					>
						Добавить
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
					<ReactSortable
						list={fields}
						setList={() => {}}
						onEnd={dropHandler}
						handle='.drag'
						{...sortableOptions}
					>
						{fields.map(item => (
							<ItemComponent key={item._id} data={item} />
						))}
					</ReactSortable>
				</FormProvider>
			</Stack>

			<DialogComponent submit={submitHandler} />
		</>
	)
}
