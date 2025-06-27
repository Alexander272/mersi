import { Table as TableContainer } from '@/components/Table/Table'
import { Head } from './Head'
import { Body } from './Body'
import { ContextMenu } from '../../modules/contextMenu/components/ContextMenuLazy'

export const Table = () => {
	return (
		<TableContainer>
			<Head />
			<Body />
			<ContextMenu />
		</TableContainer>
	)
}
