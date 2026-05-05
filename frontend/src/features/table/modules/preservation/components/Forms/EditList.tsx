import { Fragment, useEffect, useState, type FC } from 'react'
import { Collapse, Divider, List, ListItemButton, ListItemText, Stack, Typography } from '@mui/material'

import { useGetPreservationsQuery } from '../../preservationApiSlice'
import { NullDate } from '@/constants/defaultValues'
import { useGetInstrumentByIdQuery } from '@/features/table/instrumentApiSlice'
import { NoRowsOverlay } from '@/features/table/components/NoRowsOverlay/components/NoRowsOverlay'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { EditPreservationItem } from './EditItem'

type Props = {
	instrumentId: string
}

export const EditPreservationList: FC<Props> = ({ instrumentId }) => {
	const [open, setOpen] = useState('')

	const { data, isFetching } = useGetPreservationsQuery(instrumentId, { skip: !instrumentId })
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
					const title =
						'Консервация от ' +
						new Date(e.dateStart).toLocaleDateString() +
						(e.dateEnd != NullDate ? ' до ' + new Date(e.dateEnd).toLocaleDateString() : '')

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
								<EditPreservationItem data={e} instrumentId={instrumentId} />
							</Collapse>
							<Divider sx={{ width: '80%', mx: 'auto' }} />
						</Fragment>
					)
				})}
			</List>
		</>
	)
}
