import { FC, SyntheticEvent } from 'react'
import { Stack, Tab, Tabs, useTheme } from '@mui/material'

import { Fallback } from '@/components/Fallback/Fallback'
import { useGetGroupedSectionsQuery } from '../../sectionsApiSlice'

type Props = {
	item: string
	setItem: (item: string) => void
}

export const List: FC<Props> = ({ item, setItem }) => {
	const { palette } = useTheme()
	const { data, isFetching } = useGetGroupedSectionsQuery(null)

	const tabHandler = (_event: SyntheticEvent, value: string) => {
		setItem(value)
	}

	return (
		<Stack position={'relative'} minWidth={280} height={'100%'}>
			{isFetching && <Fallback position={'absolute'} zIndex={5} background={'#f5f5f557'} />}

			<Tabs
				orientation='vertical'
				value={item}
				onChange={tabHandler}
				variant='scrollable'
				sx={{
					borderRight: 1,
					borderColor: 'divider',
					maxHeight: 760,
					height: '100%',
					'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
					'.MuiTabs-scrollButtons.Mui-disabled': {
						height: 0,
					},
				}}
			>
				<Tab
					label='Добавить'
					value='new'
					sx={{
						mt: 0.5,
						textTransform: 'inherit',
						borderRadius: 3,
						transition: 'all 0.3s ease-in-out',
						maxWidth: '100%',
						minHeight: 44,
						backgroundColor: '#9ab2ef29',
						color: palette.primary.main,
						':hover': {
							backgroundColor: '#9ab2ef58',
						},
					}}
				/>

				{data?.data.map(item => {
					const arr = [
						<Tab
							key={item.id}
							label={item.title}
							value={item.id}
							disabled
							sx={{
								minHeight: 40,
								marginTop: '6px',
								padding: '8px 20px',
								position: 'relative',
								alignItems: 'flex-start',
								fontSize: '1rem',
								'&:after': {
									content: '""',
									position: 'absolute',
									left: 10,
									bottom: 0,
									width: '35%',
									height: 2,
									background: palette.primary.main,
								},
								'&.Mui-disabled': {
									color: 'inherit',
								},
							}}
						/>,
					]

					item.sections.forEach(section => {
						arr.push(
							<Tab
								key={section.id}
								label={section.name}
								value={section.id}
								sx={{
									textTransform: 'inherit',
									borderRadius: 3,
									transition: 'all 0.3s ease-in-out',
									maxWidth: '100%',
									minHeight: 48,
									':hover': {
										backgroundColor: '#f5f5f5',
									},
								}}
							/>
						)
					})
					return arr
				})}
			</Tabs>
		</Stack>
	)
}
