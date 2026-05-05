import { Fragment, useEffect, useState, type FC } from 'react'
import { Collapse, Divider, List, ListItemButton, ListItemText, Stack, Typography } from '@mui/material'

import { useGetVerificationsQuery } from '../../verificationApiSlice'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { NoRowsOverlay } from '@/features/table/components/NoRowsOverlay/components/NoRowsOverlay'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { EditVerificationItem } from './EditItem'

type Props = {
	instrumentId: string
}

export const EditVerificationList: FC<Props> = ({ instrumentId }) => {
	const [open, setOpen] = useState('')

	const { data, isFetching } = useGetVerificationsQuery(instrumentId, { skip: !instrumentId })
	const { data: item, isFetching: isFetchingItem } = useGetInstrumentByIdQuery(instrumentId, { skip: !instrumentId })

	useEffect(() => {
		if (data && data.data.length) setOpen(data.data[0].id)
	}, [data])

	const openHandler = (id: string) => () => {
		setOpen(id)
	}

	if (!data?.data.length)
		return (
			<Stack height={150} position={'relative'}>
				<NoRowsOverlay />
			</Stack>
		)

	return (
		<>
			{isFetchingItem || isFetching ? <BoxFallback /> : null}

			<Typography textAlign={'center'} fontSize={'1.2rem'} fontWeight={'bold'} sx={{ width: '100%' }}>
				{item?.data.name} ({item?.data.factoryNumber})
			</Typography>

			<List>
				{data?.data.map(e => {
					const title = `Поверка от ${new Date(e.verificationDate).toLocaleDateString()}`
					return (
						<Fragment key={e.id}>
							<ListItemButton
								onClick={openHandler(e.id)}
								selected={open === e.id}
								sx={{ borderRadius: 3 }}
							>
								<ListItemText primary={title} />
								<LeftArrowIcon
									fontSize={16}
									transform={open === e.id ? 'rotate(90deg)' : 'rotate(270deg)'}
								/>
							</ListItemButton>
							<Collapse in={open === e.id} timeout='auto'>
								<EditVerificationItem data={e} />
							</Collapse>
							<Divider sx={{ width: '80%', mx: 'auto' }} />
						</Fragment>
					)
				})}
			</List>
		</>
	)
}
