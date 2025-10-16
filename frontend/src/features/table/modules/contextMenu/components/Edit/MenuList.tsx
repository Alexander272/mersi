import { FC } from 'react'
import { IconButton, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material'

import { useAppDispatch } from '@/hooks/redux'
import { useGetContextMenuQuery } from '../../contextMenuApiSlice'
import { changeDialogIsOpen } from '@/features/dialog/dialogSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { EditIcon } from '@/components/Icons/EditIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { EditContextDialog } from './Dialog'
import { IContextMenu } from '../../types/context'

type Props = {
	section: string
}

export const ContextMenuList: FC<Props> = ({ section }) => {
	const dispatch = useAppDispatch()

	const { data, isFetching } = useGetContextMenuQuery(
		{ section: section, isFull: true },
		{ skip: !section || section == 'new' }
	)

	const addHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'EditContextMenu', isOpen: true, context: { section } }))
	}

	const editHandler = (data: IContextMenu) => () => {
		dispatch(changeDialogIsOpen({ variant: 'EditContextMenu', isOpen: true, context: { data, section } }))
	}

	return (
		<TableContainer sx={{ position: 'relative' }}>
			{isFetching && <BoxFallback />}

			<Table size='small'>
				<TableHead>
					<TableRow>
						<TableCell>№ в списке</TableCell>
						<TableCell>Код пункта</TableCell>
						<TableCell>Название (необязательно)</TableCell>
						<TableCell>Правило для отображения</TableCell>
						<TableCell>
							<IconButton onClick={addHandler} size='large' disabled={section == 'new'}>
								<PlusIcon fontSize={12} fill={section == 'new' ? '#bdbdbd' : '#05287f'} />
							</IconButton>
						</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{data?.data.map(d => (
						<TableRow key={d.id}>
							<TableCell>{d.position}</TableCell>
							<TableCell>{d.name}</TableCell>
							<TableCell>{d.label}</TableCell>
							<TableCell>{d.rule}</TableCell>
							<TableCell>
								<IconButton onClick={editHandler(d)} size='large'>
									<EditIcon fontSize={12} fill={'#6e6e6e'} />
								</IconButton>
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>

			<EditContextDialog />
		</TableContainer>
	)
}
