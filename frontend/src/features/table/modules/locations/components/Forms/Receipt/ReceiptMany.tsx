import { FC, useEffect, useState } from 'react'
import {
	Button,
	Divider,
	FormControl,
	InputLabel,
	MenuItem,
	OutlinedInput,
	Select,
	SelectChangeEvent,
	Stack,
	Theme,
	Typography,
	useTheme,
} from '@mui/material'
import { IFilter } from '@/features/table/types/params'
import { useCheckPermission } from '@/features/user/hooks/check'
import { PermRules } from '@/constants/permissions'
import { useGetDepartmentsQuery } from '@/features/departments/departmentApiSlice'
import { useGetResponsibleByUserQuery } from '@/features/departments/responsibleApiSlice'
import { useGetSentSIQuery } from '@/features/table/siApiSlice'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'
import type { ISentSI } from '@/features/table/types/si'
import { useReceivingMutation } from '../../../locationsApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { IFetchError } from '@/app/types/error'
import { toast } from 'react-toastify'
import { Fallback } from '@/components/Fallback/Fallback'
import { CheckboxGroup } from '@/components/CheckboxGroup/CheckboxGroup'

const defFilter: IFilter[] = [
	{
		field: 'department',
		fieldType: 'list',
		compareType: 'null',
		value: '',
	},
]

export const ReceiptMany: FC = () => {
	const theme = useTheme()
	const section = useAppSelector(getSection)
	const hasResWrite = useCheckPermission(PermRules.Reserve.Write)

	const [enable, setEnable] = useState<string[]>([])
	const [filters, setFilters] = useState(defFilter)
	const { data: departments, isFetching: isFetchDepartments } = useGetDepartmentsQuery(section?.realmId || '', {
		skip: !section?.realmId,
	})
	// const { data: departmentsByUser } = useGetDepartmentsByUserQuery(null)
	const { data: departmentsByUser } = useGetResponsibleByUserQuery(null)

	const { data, isFetching } = useGetSentSIQuery({ section: section?.id || '', filters }, { skip: !section?.id })

	useEffect(() => {
		if (departmentsByUser) setEnable(departmentsByUser.data.map(d => d.departmentId))
	}, [departmentsByUser])
	useEffect(() => {
		if (hasResWrite) setFilters([{ ...defFilter[0], compareType: 'in', value: enable.join(',') }])
	}, [hasResWrite, enable])

	const handleChange = (event: SelectChangeEvent<typeof enable>) => {
		const { value } = event.target
		setEnable(typeof value === 'string' ? value.split(',') : value)
		const newFilter = {
			...defFilter[0],
			compareType: 'in' as const,
			value: typeof value === 'string' ? value : value.join(','),
		}
		setFilters([newFilter])
	}

	return (
		<Stack position={'relative'} mt={-2.5}>
			{hasResWrite ? (
				<FormControl sx={{ m: 1, width: '100%', mb: 2 }}>
					<InputLabel id='departments'>Подразделение</InputLabel>
					<Select
						labelId='departments'
						multiple
						value={enable}
						onChange={handleChange}
						input={<OutlinedInput label='Подразделение' />}
						// MenuProps={MenuProps}
					>
						{departments?.data.map(d => (
							<MenuItem key={d.id} value={d.id} style={getStyles(d.id, enable, theme)}>
								{d.name}
							</MenuItem>
						))}
					</Select>
				</FormControl>
			) : null}

			{isFetchDepartments || isFetching ? <Fallback /> : <GroupedList data={data?.data || []} />}
		</Stack>
	)
}

type GroupedListProps = {
	data: ISentSI[]
}
const GroupedList: FC<GroupedListProps> = ({ data }) => {
	const dispatch = useAppDispatch()

	const defState = data.flatMap(item => item.si).reduce((a, v) => a.set(v.id, true), new Map<string, boolean>())
	const [checked, setChecked] = useState<Map<string, boolean>>(defState)
	const hasResWrite = useCheckPermission(PermRules.Reserve.Write)

	const [receiving, { isLoading }] = useReceivingMutation()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Receive', isOpen: false }))
	}

	const receiveHandler = async () => {
		const value: string[] = []
		checked.forEach((v, k) => {
			if (v) value.push(k)
		})

		const payload = {
			instrumentId: value,
			status: hasResWrite ? 'used' : 'reserve',
		}
		console.log(payload)

		try {
			await receiving(payload).unwrap()
			closeHandler()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
			console.log(error)
		}
	}

	if (!data.length) return <Typography my={2}>Инструменты для получения не найдены</Typography>
	return (
		<>
			{isLoading && (
				<Fallback
					position={'absolute'}
					top={'50%'}
					left={'50%'}
					transform={'translate(-50%, -50%)'}
					height={160}
					width={160}
					borderRadius={3}
					zIndex={15}
					backgroundColor={'#fafafa'}
				/>
			)}

			<Stack width={'100%'}>
				{data.map(item => (
					<CheckboxGroup
						key={item.place}
						checked={checked}
						data={{
							name: item.place,
							list: item.si.map(si => ({ ...si, name: `${si.name} (${si.factoryNumber})` })) || [],
						}}
						onChange={setChecked}
					/>
				))}
			</Stack>

			<Divider sx={{ width: '50%', alignSelf: 'center', mt: 1 }} />
			<Stack spacing={2} direction={'row'} mt={2} justifyContent={'center'}>
				<Button onClick={receiveHandler} variant='outlined' sx={{ width: '50%' }}>
					Подтвердить
				</Button>
			</Stack>
		</>
	)
}

function getStyles(name: string, array: string[], theme: Theme) {
	return {
		fontWeight: array.indexOf(name) === -1 ? theme.typography.fontWeightRegular : theme.typography.fontWeightMedium,
	}
}
