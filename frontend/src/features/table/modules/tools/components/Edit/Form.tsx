import { FC } from 'react'
import {
	Button,
	Checkbox,
	Divider,
	FormControl,
	FormControlLabel,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	TextField,
} from '@mui/material'
import { Controller, useForm } from 'react-hook-form'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { IToolsMenu, IToolsMenuDTO } from '../../types/toolsMenu'
import { useAppDispatch } from '@/hooks/redux'
import { useGetRulesQuery } from '@/features/rules/rulesApiSlice'
import { useCreateToolsMenuMutation, useUpdateToolsMenuMutation } from '../../toolsMenuApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { Fallback } from '@/components/Fallback/Fallback'

const defaultData: IToolsMenuDTO = {
	id: '',
	sectionId: '',
	position: 1,
	name: '',
	label: '',
	rule: '',
	ruleItemId: '',
	canBeFavorite: false,
}

const ToolCodes = [
	'graphic',
	'export',
	'verification',
	'send-to-verification',
	'locations',
	'reserve',
	'receive',
	'repair-info',
	'preservation-info',
	'transfer-to-save',
	'transfer-to-department',
	'write-off',
	'departments',
	'places',
]

type Props = {
	data?: IToolsMenu
	section: string
}

export const EditToolsForm: FC<Props> = ({ data, section }) => {
	const dispatch = useAppDispatch()

	const {
		control,
		handleSubmit,
		formState: { dirtyFields },
	} = useForm<IToolsMenuDTO>({ defaultValues: data || defaultData })

	const { data: rules, isFetching } = useGetRulesQuery(null)
	const [create, { isLoading: creating }] = useCreateToolsMenuMutation()
	const [update, { isLoading: updating }] = useUpdateToolsMenuMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditToolsMenu', isOpen: false }))
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
				toast.success('Элемент меню "Инструменты" создан')
			} else {
				await update(form).unwrap()
				toast.success('Элемент меню "Инструменты" обновлен')
			}
			dispatch(changeDialogIsOpen({ variant: 'EditToolsMenu', isOpen: false }))
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
				<InputLabel id={'name'}>Код пункта</InputLabel>

				<Controller
					control={control}
					name='name'
					rules={{ required: true }}
					// render={({ field }) => <TextField {...field} label={'Код пункта'} fullWidth />}
					render={({ field }) => (
						<Select {...field} labelId='name' label={'Код пункта'} fullWidth>
							{ToolCodes.map(code => (
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

			<Controller
				control={control}
				name='canBeFavorite'
				render={({ field }) => (
					<FormControl>
						<FormControlLabel control={<Checkbox {...field} />} label='Можно добавить в контекстное меню' />
					</FormControl>
				)}
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
