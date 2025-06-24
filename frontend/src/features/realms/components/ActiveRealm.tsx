import { useEffect } from 'react'
import { MenuItem, Select, SelectChangeEvent, useTheme } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { setPermissions, setRole } from '@/features/user/userSlice'
import { setPage } from '@/features/table/tableSlice'
import { useChooseRealmMutation, useGetRealmsByUserQuery } from '../realmsApiSlice'
import { getRealm, setRealm } from '../realmSlice'

export const ActiveRealm = () => {
	const { palette } = useTheme()

	const realm = useAppSelector(getRealm)
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetRealmsByUserQuery(null)
	const [choose] = useChooseRealmMutation()

	useEffect(() => {
		if (!data) return
		const founded = data.data.find(e => e.id === realm?.id)
		if (founded) return
		dispatch(setRealm(data.data[0]))
	}, [data, dispatch, realm])

	const changeHandler = async (event: SelectChangeEvent) => {
		const value = data?.data.find(e => e.id === event.target.value)
		if (!value) return

		dispatch(setRealm(value))
		try {
			const payload = await choose(value.id).unwrap()
			dispatch(setRole(payload.data.role))
			dispatch(setPermissions(payload.data.permissions))
			// dispatch(setFilters(payload.data.filters))
			dispatch(setPage(1))
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	if ((data?.data.length || 0) < 2) return null
	return (
		<Select
			value={realm?.id || ''}
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
