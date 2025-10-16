import { useState } from 'react'
import { Stack } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { ToolsMenuList } from '@/features/table/modules/tools/components/Edit/ToolsList'

export default function Tools() {
	const [item, setItem] = useState<string>('new')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'}>
			<List item={item} setItem={itemHandler} readonly />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
				<ToolsMenuList section={item} />
			</Stack>
		</Stack>
	)
}
