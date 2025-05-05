import { Box, Stack, Typography } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { Pagination } from '@/components/Pagination/Pagination'
import { getSelected, getTablePage, getTableSize, setPage } from '../../tableSlice'

export const Footer = () => {
	const size = useAppSelector(getTableSize)
	const page = useAppSelector(getTablePage)
	const selected = useAppSelector(getSelected)

	const dispatch = useAppDispatch()

	// const { data } = useGetAllData()

	const setPageHandler = (page: number) => {
		dispatch(setPage(page))
	}

	// const totalPages = Math.ceil((data?.total || 1) / size)
	const totalPages = 10

	return (
		<Box display={'grid'} alignItems={'center'} gridTemplateColumns={'repeat(3, 1fr)'} mt={1} mx={2}>
			<Typography pr={1.5} mr={'auto'}>
				Строк выбрано: {Object.keys(selected).length}
			</Typography>

			{totalPages > 1 ? (
				<Pagination page={page} totalPages={totalPages} onClick={setPageHandler} sx={{ marginX: 'auto' }} />
			) : (
				<span />
			)}

			{/* {data?.data.length ? (
				<Stack direction={'row'} alignItems={'center'} justifyContent={'flex-end'}>
					<TableSize total={data?.total || 1} />
					<Typography sx={{ ml: 2 }}>
						{(page - 1) * size || 1}-{(page - 1) * size + (data?.data.length || 0)} из {data?.total}
					</Typography>
				</Stack>
			) : null} */}
		</Box>
	)
}
