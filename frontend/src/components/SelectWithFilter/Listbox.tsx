import { Box, Divider, Stack, useTheme } from '@mui/material'
import { createContext, forwardRef, useContext, useRef } from 'react'
import { FixedSizeList, ListChildComponentProps } from 'react-window'

import { CheckIcon } from '../Icons/CheckIcon'

export const ListboxComponent = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLElement>>(function ListboxComponent(
	props,
	ref
) {
	const { children, ...other } = props
	const listRef = useRef<FixedSizeList>(null)
	const itemData: React.ReactElement<unknown>[] = []
	;(children as React.ReactElement<unknown>[]).forEach(
		(
			item: React.ReactElement<unknown> & {
				children?: React.ReactElement<unknown>[]
			}
		) => {
			itemData.push(item)
			itemData.push(...(item.children || []))
		}
	)

	const itemCount = itemData.length
	const itemSize = 39

	const getChildSize = () => {
		return itemSize
	}

	const getHeight = () => {
		if (itemCount > 8) {
			return 8 * itemSize
		}
		return itemData.map(getChildSize).reduce((a, b) => a + b + 2, 0)
	}

	return (
		<div ref={ref}>
			<OuterElementContext.Provider value={other}>
				<FixedSizeList
					itemData={itemData}
					height={getHeight() + 24}
					width='100%'
					ref={listRef}
					outerElementType={OuterElementType}
					innerElementType='ul'
					itemSize={itemSize}
					overscanCount={6}
					itemCount={itemCount}
				>
					{Row}
				</FixedSizeList>
			</OuterElementContext.Provider>
		</div>
	)
})

const Row = (props: ListChildComponentProps) => {
	const { palette } = useTheme()

	const { data, index, style } = props
	const dataSet = data[index]
	const inlineStyle = {
		...style,
		width: `calc(${style.width} - 12px)`, // style.width as number,
		top: style.top as number,
		// left: (style.left as number) + 6,
	}

	const { key, ...optionProps } = dataSet[0]
	const option = dataSet[1]
	const { selected } = dataSet[2]

	return (
		<li key={key} style={inlineStyle}>
			<Stack width={'100%'} {...optionProps}>
				<Stack direction={'row'} width={'100%'}>
					<Box
						display={'flex'}
						justifyContent={'center'}
						alignItems={'center'}
						sx={{
							width: 20,
							minWidth: 20,
							height: 20,
							mr: 2,
							ml: 1,
							borderRadius: 1,
							border: '1px solid #afafaf',
						}}
					>
						{selected ? <CheckIcon fontSize={14} fill={palette.primary.main} /> : null}
					</Box>
					<Box
						sx={{
							overflow: 'hidden',
							textOverflow: 'ellipsis',
							whiteSpace: 'nowrap',
						}}
					>
						{option.name}
					</Box>
				</Stack>
				<Divider flexItem sx={{ mt: 1 }} />
			</Stack>
		</li>
	)
}

const OuterElementContext = createContext({})

const OuterElementType = forwardRef<HTMLDivElement>((props, ref) => {
	const outerProps = useContext(OuterElementContext)
	return <div ref={ref} {...props} {...outerProps} />
})
