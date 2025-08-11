import { useCallback, useEffect } from 'react'
import { Box } from '@mui/material'

import { PermRules } from '@/constants/permissions'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useLazyExportQuery } from '../exportApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { getFilters, getSearch, getSort, getStatus } from '@/features/table/tableSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

export const ExportDialog = () => {
	const status = useAppSelector(getStatus)
	const section = useAppSelector(getSection)
	const search = useAppSelector(getSearch)
	const sort = useAppSelector(getSort)
	const filters = useAppSelector(getFilters)

	const all = useCheckPermission(PermRules.Location.Write)

	const modal = useAppSelector(getDialogState('Export'))
	const dispatch = useAppDispatch()

	const [exportData] = useLazyExportQuery()

	const exportFunc = useCallback(async () => {
		await exportData({ section: section?.id || '', status, all, sort, filters, search })
		dispatch(changeDialogIsOpen({ variant: 'Export', isOpen: false }))
	}, [all, exportData, filters, search, section?.id, sort, status, dispatch])

	useEffect(() => {
		if (modal?.isOpen) exportFunc()
	}, [exportFunc, modal?.isOpen])

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Export', isOpen: false }))
	}

	return (
		<Dialog
			title={'Экспорт данных'}
			body={<Body />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='xs'
			fullWidth
		/>
	)
}

const Body = () => {
	return (
		<Box position={'relative'} mt={-3} width={'100%'} height={180}>
			<BoxFallback />
		</Box>
	)
}
