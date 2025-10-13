import { MenuItem, Select, SelectChangeEvent, useTheme } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useGetStatusesQuery } from '../../statusApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getStatus, setStatus } from '@/features/table/tableSlice'

export const Status = () => {
	const { palette } = useTheme()

	const status = useAppSelector(getStatus)
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetStatusesQuery(section?.id || '', { skip: !section?.id })

	const changeHandler = (event: SelectChangeEvent) => {
		const value = event.target.value
		dispatch(setStatus(value))
	}

	return (
		<Select
			value={status}
			onChange={changeHandler}
			disabled={isFetching}
			sx={{
				color: palette.primary.main,
				fontSize: '1.2rem',
				boxShadow: 'none',
				'.MuiOutlinedInput-notchedOutline': { border: 0 },
				'&.MuiOutlinedInput-root:hover .MuiOutlinedInput-notchedOutline': {
					border: 0,
				},
				'&.MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline': {
					border: 0,
				},
				'.MuiOutlinedInput-input': { padding: '6.5px 10px' },
			}}
		>
			{data?.data.map(s => (
				<MenuItem key={s.id} value={s.value}>
					{s.label}
				</MenuItem>
			))}

			{/* <MenuItem value={'work'}>Основные</MenuItem>
			<MenuItem value={'repair'}>На ремонте</MenuItem>
			{section?.id == '46ba9e17-65c7-474b-8c47-7975ab4319d5' && [
				<MenuItem value={'archived'}>Законсервированные</MenuItem>,
				<MenuItem value={'saved'}>На хранении</MenuItem>,
				<MenuItem value={'transferred'}>Переданные</MenuItem>,
			]}
			<MenuItem value={'decommissioning'}>Непригодные</MenuItem> */}
		</Select>
	)
}
