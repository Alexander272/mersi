import { useState } from 'react'
import { Stack } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { ContextMenuList } from '@/features/table/modules/contextMenu/components/Edit/MenuList'

export default function ContextMenu() {
	const [item, setItem] = useState<string>('new')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'}>
			<List item={item} setItem={itemHandler} readonly />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
				<ContextMenuList section={item} />
			</Stack>
		</Stack>
	)
}
