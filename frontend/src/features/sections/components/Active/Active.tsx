import { useEffect } from 'react'
import { MenuItem, Select, SelectChangeEvent, useTheme } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getRealm } from '@/features/realms/realmSlice'
import { useGetSectionsQuery } from '../../sectionsApiSlice'
import { getSection, setSection } from '../../sectionSlice'

export const ActiveSection = () => {
	const { palette } = useTheme()

	const realm = useAppSelector(getRealm)
	const section = useAppSelector(getSection)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetSectionsQuery(realm?.id || '', { skip: !realm?.id })

	useEffect(() => {
		if (!data) return
		const founded = data.data.find(e => e.id === realm?.id)
		if (founded) return
		dispatch(setSection(data.data[0]))
	}, [data, dispatch, realm])

	const changeHandler = (event: SelectChangeEvent) => {
		const value = data?.data.find(e => e.id === event.target.value)
		if (!value) return

		dispatch(setSection(value))
	}

	return (
		<Select
			value={section?.id || ''}
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
			<MenuItem value='' disabled>
				Выберите область
			</MenuItem>
			{data?.data.map(item => (
				<MenuItem key={item.id} value={item.id}>
					{item.name}
				</MenuItem>
			))}
		</Select>
	)
}
