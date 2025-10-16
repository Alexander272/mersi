import { FC } from 'react'
import { Button, Divider, FormControl, InputLabel, MenuItem, Select, Stack, TextField } from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IContextMenu, IContextMenuDTO } from '../../types/context'
import { useAppDispatch } from '@/hooks/redux'
import { useGetRulesQuery } from '@/features/rules/rulesApiSlice'
import { useCreateContextMenuMutation, useUpdateContextMenuMutation } from '../../contextMenuApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Fallback } from '@/components/Fallback/Fallback'

const defaultData: IContextMenuDTO = {
	id: '',
	sectionId: '',
	position: 1,
	name: '',
	label: '',
	rule: '',
	ruleItemId: '',
}

const ContextCodes = [
	'create',
	'edit',
	'change-position',
	'verification',
	'sent-to-verification',
	'location',
	'reserve',
	'receive',
	'forced',
	'cancel',
	'repair-info',
	'preservation-info',
	'transfer-to-save',
	'transfer-to-department',
	'write-off',
	'history',
]

type Props = {
	data?: IContextMenu
	section: string
}

export const ExitContextForm: FC<Props> = ({ data, section }) => {
	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IContextMenuDTO>({ defaultValues: data || defaultData })

	const { data: rules, isFetching } = useGetRulesQuery(null)
	const [create, { isLoading: creating }] = useCreateContextMenuMutation()
	const [update, { isLoading: updating }] = useUpdateContextMenuMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditContextMenu', isOpen: false }))
	}

	const saveHandler = handleSubmit(async form => {
		console.log('save', form, dirtyFields)

		form.position = +form.position
		form.sectionId = section
		form.name = form.name.trim()
		form.label = form.label.trim()
		form.ruleItemId = rules?.data.find(r => `${r.name}:${r.method}` === form.rule)?.id || ''

		try {
			if (!data?.id) {
				await create(form).unwrap()
				toast.success('Элемент контекстного меню создан')
			} else {
				await update(form).unwrap()
				toast.success('Элемент контекстного меню обновлен')
			}
			dispatch(changeDialogIsOpen({ variant: 'EditContextMenu', isOpen: false }))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	return (
		<Stack component={'form'} position={'relative'} spacing={2} onSubmit={saveHandler} mt={-2}>
			<Controller
				control={control}
				name='position'
				rules={{ required: true }}
				render={({ field }) => (
					<TextField {...field} label={'№ в списке'} fullWidth slotProps={{ input: { type: 'number' } }} />
				)}
			/>

			<FormControl>
				<InputLabel id={'name'}>Правило для отображения</InputLabel>
				<Controller
					control={control}
					name='name'
					rules={{ required: true }}
					// render={({ field }) => <TextField {...field} label={'Код пункта'} fullWidth />}
					render={({ field }) => (
						<Select {...field} labelId='name' label={'Код пункта'} fullWidth>
							{ContextCodes.map(code => (
								<MenuItem key={code} value={code}>
									{code}
								</MenuItem>
							))}
						</Select>
					)}
				/>
			</FormControl>

			<Controller
				control={control}
				name='label'
				render={({ field }) => <TextField {...field} label={'Название (необязательно)'} fullWidth />}
			/>

			<FormControl>
				<InputLabel id={'rule'}>Правило для отображения</InputLabel>

				<Controller
					control={control}
					name='rule'
					rules={{ required: true }}
					render={({ field }) => (
						<Select {...field} labelId='rule' label={'Правило для отображения'} fullWidth>
							{rules?.data.map(rule => (
								<MenuItem key={rule.id} value={`${rule.name}:${rule.method}`}>
									{rule.name}:{rule.method} ({rule.description})
								</MenuItem>
							))}
						</Select>
					)}
				/>
			</FormControl>

			<Divider sx={{ width: '50%', alignSelf: 'center' }} />
			<Stack spacing={2} direction={'row'}>
				<Button type='submit' variant='contained' fullWidth>
					Сохранить
				</Button>
				<Button onClick={closeHandler} variant='outlined' fullWidth>
					Отмена
				</Button>
			</Stack>

			{isFetching || creating || updating ? (
				<Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} mt={'0!important'} />
			) : null}
		</Stack>
	)
}
