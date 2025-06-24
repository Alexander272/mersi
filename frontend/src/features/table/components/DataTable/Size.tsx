import { FC } from 'react'
import { MenuItem, Select, SelectChangeEvent, Stack, Typography } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getTablePage, getTableSize, setPage, setSize } from '../../tableSlice'

const limits = [15, 30, 50, 100]

type Props = {
	total: number
}

export const Size: FC<Props> = ({ total }) => {
	const size = useAppSelector(getTableSize)
	const page = useAppSelector(getTablePage)

	const dispatch = useAppDispatch()

	const changeSize = (event: SelectChangeEvent<number>) => {
		const newSize = +event.target.value

		const maxPage = Math.ceil(total / newSize)
		if (page > maxPage) dispatch(setPage(maxPage))

		dispatch(setSize(newSize))
	}

	return (
		<Stack direction={'row'} alignItems={'center'}>
			<Typography>Строк на странице:</Typography>

			<Select
				name='table_size'
				value={size}
				onChange={changeSize}
				sx={{
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
				{limits.map(l => (
					<MenuItem key={l} value={l}>
						{l}
					</MenuItem>
				))}
			</Select>
		</Stack>
	)
}
