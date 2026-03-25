import { useCallback, useEffect } from 'react'
import { Box } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { useLazyMakeAccountingLogQuery } from '../exportApiSlice'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

export const AccountingLogDialog = () => {
	const section = useAppSelector(getSection)

	const modal = useAppSelector(getDialogState('Log'))
	const dispatch = useAppDispatch()

	const [exportData] = useLazyMakeAccountingLogQuery()

	const exportFunc = useCallback(async () => {
		await exportData({ section: section?.id || '', gte: '', lte: '' })
		dispatch(changeDialogIsOpen({ variant: 'Log', isOpen: false }))
	}, [exportData, section?.id, dispatch])

	useEffect(() => {
		if (modal?.isOpen) exportFunc()
	}, [exportFunc, modal?.isOpen])

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'Log', isOpen: false }))
	}

	return (
		<Dialog
			title={'Формирование журнала учета'}
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
