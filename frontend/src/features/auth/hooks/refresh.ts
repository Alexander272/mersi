import { useEffect, useState } from 'react'

import { useAppDispatch } from '@/hooks/redux'
import { setUser } from '@/features/user/userSlice'
import { useRefreshQuery } from '../authApiSlice'

export function useRefresh() {
	const [ready, setReady] = useState(false)

	const { data, isError, isSuccess } = useRefreshQuery(null)

	const dispatch = useAppDispatch()

	useEffect(() => {
		if (isSuccess) {
			dispatch(setUser(data.data))
			setReady(true)
		}
	}, [isSuccess, data, dispatch])

	useEffect(() => {
		if (isError) setReady(true)
	}, [isError])

	return { ready }
}

// export function useRefresh() {
// 	const [ready, setReady] = useState(false)

// 	const { data, isError, isSuccess } = useRefreshQuery(null)

// 	const dispatch = useAppDispatch()

// 	useEffect(() => {
// 		if (!isError && data) dispatch(setUser(data.data))
// 		setReady(true)
// 		console.log('refresh', isError, data)
// 		console.log('ready', ready)
// 	}, [isError, data, dispatch])

// 	return { ready }
// }
