import { Table as TableContainer } from '@/components/Table/Table'
import { Head } from './Head'
import { Body } from './Body'

export const Table = () => {
	return (
		<TableContainer>
			<Head />
			<Body />
			{/* <ContextMenu /> */}
		</TableContainer>
	)
}
