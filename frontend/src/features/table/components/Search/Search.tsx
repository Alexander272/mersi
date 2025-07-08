import { ChangeEvent } from 'react'
import { InputAdornment, TextField } from '@mui/material'

import { useAppDispatch } from '@/hooks/redux'
import { useDebounceFunc } from '@/hooks/useDebounceFunc'
import { SearchIcon } from '@/components/Icons/SearchIcon'
import { setSearch } from '../../tableSlice'

export const Search = () => {
	// const search = useAppSelector(getSearch)
	const dispatch = useAppDispatch()

	const searchHandler = useDebounceFunc(v => {
		dispatch(setSearch(v as string))
	}, 700)

	const changeValueHandler = (event: ChangeEvent<HTMLInputElement>) => {
		// setValue(event.target.value)
		searchHandler(event.target.value)
	}

	return (
		<TextField
			// value={search.value}
			name='search'
			onChange={changeValueHandler}
			slotProps={{
				input: {
					startAdornment: (
						<InputAdornment position='start'>
							<SearchIcon fontSize={16} />
						</InputAdornment>
					),
					// endAdornment: (
					// 	<InputAdornment position='end'>
					// 		<Setting />
					// 	</InputAdornment>
					// ),
				},
			}}
			placeholder='Поиск...'
			sx={{ width: 350 }}
		/>
	)
}
