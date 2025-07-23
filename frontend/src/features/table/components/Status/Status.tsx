import { MenuItem, Select, SelectChangeEvent, useTheme } from '@mui/material'

import type { Status as StatusType } from '../../types/si'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getSection } from '@/features/sections/sectionSlice'
import { getStatus, setStatus } from '../../tableSlice'

export const Status = () => {
	const { palette } = useTheme()

	const status = useAppSelector(getStatus)
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const changeHandler = (event: SelectChangeEvent) => {
		const value = event.target.value as StatusType
		dispatch(setStatus(value))
	}

	return (
		<Select
			value={status}
			onChange={changeHandler}
			// disabled={isFetching}
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
			<MenuItem value={'work'}>Основные</MenuItem>
			<MenuItem value={'repair'}>На ремонте</MenuItem>
			<MenuItem value={'decommissioning'}>Непригодные</MenuItem>
			{section?.id == '46ba9e17-65c7-474b-8c47-7975ab4319d5' && (
				<MenuItem value={'transferred'}>Переданные</MenuItem>
			)}
		</Select>
	)
}
