import { FC, useEffect, useMemo, useState } from 'react'
import { Box, IconButton, Tab, Tabs } from '@mui/material'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { EditRepairList } from '@/features/table/modules/repair/components/Form/EditList'
import { EditTransferToDepList } from '@/features/table/modules/transferToDep/components/Forms/EditList'
import { EditTransferToSaveList } from '@/features/table/modules/transferToSave/components/Forms/EditList'
import { EditWriteOffList } from '@/features/table/modules/writeOff/components/Forms/EditList'
import { EditPreservationList } from '@/features/table/modules/preservation/components/Forms/EditList'
import { EditVerificationList } from '@/features/table/modules/verification/components/Forms/EditList'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { getSection } from '@/features/sections/sectionSlice'
import { useGetToolsMenuQuery } from '../../modules/tools/toolsMenuApiSlice'
import { Fallback } from '@/components/Fallback/Fallback'

type Context = { id: string }

export const UpdateDetailsDialog = () => {
	const modal = useAppSelector(getDialogState('UpdateTableDetails'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'UpdateTableDetails', isOpen: false }))
	}

	const context = modal?.context as Context | undefined

	return (
		<Dialog
			title={'Редактировать'}
			headerActions={
				<IconButton onClick={closeHandler} size='large' sx={{ fill: '#505050', mr: 2 }}>
					<TimesIcon fontSize={12} />
				</IconButton>
			}
			body={<Form id={context?.id || ''} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}

const tabStyle = {
	textTransform: 'inherit',
	borderRadius: 3,
	transition: 'all 0.3s ease-in-out',
	minHeight: 48,
	flexGrow: 1,
	flexShrink: 1,
	maxWidth: '100%',
	width: '100%',
	':hover': {
		backgroundColor: '#f5f5f5',
	},
}

const TAB_CONFIG = [
	{ value: 'verification', label: 'Поверки', component: EditVerificationList },
	{ value: 'repair-info', label: 'Ремонт', component: EditRepairList },
	{ value: 'preservation-info', label: 'Консервация', component: EditPreservationList },
	{ value: 'transfer-to-department', label: 'Перемещение в подразделение', component: EditTransferToDepList },
	{ value: 'transfer-to-save', label: 'Перемещение на хранение', component: EditTransferToSaveList },
	{ value: 'write-off', label: 'Списание', component: EditWriteOffList },
] as const

const Form: FC<{ id: string }> = ({ id }) => {
	const [value, setValue] = useState('repair')

	const section = useAppSelector(getSection)

	const { data, isFetching } = useGetToolsMenuQuery({ section: section?.id || '' }, { skip: !section?.id })

	const allowedTabs = useMemo(() => {
		if (!data?.data) return []

		// Создаем Set для быстрой проверки прав
		const allowedRules = new Set(data.data.map(d => d.name))

		return TAB_CONFIG.filter(tab => allowedRules.has(tab.value))
	}, [data])

	useEffect(() => {
		if (allowedTabs.length) setValue(allowedTabs[0].value)
	}, [allowedTabs])

	const tabContent = {
		verification: <EditVerificationList instrumentId={id} />,
		'repair-info': <EditRepairList instrumentId={id} />,
		'preservation-info': <EditPreservationList instrumentId={id} />,
		'transfer-to-department': <EditTransferToDepList instrumentId={id} />,
		'transfer-to-save': <EditTransferToSaveList instrumentId={id} />,
		'write-off': <EditWriteOffList instrumentId={id} />,
	}

	const tabHandler = (_event: React.SyntheticEvent, newValue: string) => {
		setValue(newValue)
	}

	if (isFetching) return <Fallback />
	return (
		<Box position={'relative'} mt={-2}>
			<Tabs
				value={value}
				onChange={tabHandler}
				variant='scrollable'
				sx={{
					mt: 1,
					mb: 2,
					borderBottom: 1,
					borderColor: '#00000014',
					'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
					'.MuiTabs-scrollButtons.Mui-disabled': {
						height: 0,
					},
				}}
			>
				{allowedTabs.map(t => (
					<Tab key={t.value} label={t.label} value={t.value} sx={tabStyle} />
				))}
			</Tabs>

			{tabContent[value as keyof typeof tabContent]}
		</Box>
	)
}
