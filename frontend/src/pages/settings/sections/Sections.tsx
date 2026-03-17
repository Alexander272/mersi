import { useState } from 'react'
import { Stack } from '@mui/material'

import { List } from '@/features/sections/components/List/List'
import { Form } from '@/features/sections/components/Form/Form'
import { Columns } from '@/features/sections/modules/columns/components/Columns/Columns'
import { ChildrenTabs } from '@/features/sections/components/Tabs/ChildrenTabs'
import { FieldsList } from '@/features/sections/modules/form/components/List/FieldsList'

export default function Sections() {
	const [item, setItem] = useState('new')
	const [child, setChild] = useState('columns')

	const itemHandler = (data: string) => {
		setItem(data)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'} pl={2}>
			<List item={item} setItem={itemHandler} />
			<Stack width={'100%'} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1, pr: 2 }}>
				<Form section={item} setSection={itemHandler} />
				{item != 'new' && (
					<>
						<ChildrenTabs value={child} onChange={setChild} />
						{child == 'columns' && <Columns section={item} />}
						{child == 'form' && <FieldsList section={item} />}
					</>
				)}
			</Stack>
		</Stack>
	)
}
