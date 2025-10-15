import { useState } from 'react'
import { Stack } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { StatusesList } from '@/features/table/modules/status/components/Edit/StatusesList'

export default function Statuses() {
	const [item, setItem] = useState<string>('new')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'}>
			<List item={item} setItem={itemHandler} readonly />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
				<StatusesList section={item} />
			</Stack>
		</Stack>
	)
}
