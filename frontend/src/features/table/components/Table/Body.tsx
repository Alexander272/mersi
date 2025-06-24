import { FixedSizeList } from 'react-window'

import { MaxSize, RowHeight, Size } from '@/constants/defaultValues'
import { useAppSelector } from '@/hooks/redux'
import { useGetColumnsQuery } from '@/features/sections/modules/columns/columnsApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { TableBody } from '@/components/Table/TableBody'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { useGetSI } from '../../hooks/getSI'
import { getTableSize } from '../../tableSlice'
import { useCalcWidth } from '../../utils/calcWidth'
import { NoRowsOverlay } from '../NoRowsOverlay/components/NoRowsOverlay'
import { Row } from './Row'

export const Body = () => {
	const section = useAppSelector(getSection)
	const size = useAppSelector(getTableSize)

	const { data, isFetching, isLoading } = useGetSI()
	const { data: columns } = useGetColumnsQuery(section?.id || '', { skip: !section?.id })

	const { width } = useCalcWidth(columns?.data || [])

	if (!isLoading && !data?.total) return <NoRowsOverlay />
	return (
		<TableBody>
			{isFetching || isLoading ? <BoxFallback /> : null}

			{data && (
				<FixedSizeList
					overscanCount={10}
					height={RowHeight * (size > Size ? MaxSize : Size)}
					itemCount={data.data.length > (size || Size) ? size || Size : data.data.length}
					itemSize={RowHeight}
					width={width}
				>
					{({ index, style }) => <Row item={data.data[index]} sx={style} />}
				</FixedSizeList>
			)}
		</TableBody>
	)
}
