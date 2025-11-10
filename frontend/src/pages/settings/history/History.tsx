import { useState } from 'react'
import { Stack } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { HistoryList } from '@/features/table/modules/history/components/Edit/List'

export default function History() {
	const [item, setItem] = useState<string>('new')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'}>
			<List item={item} setItem={itemHandler} readonly />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
				<HistoryList section={item} />
			</Stack>
		</Stack>
	)
}
