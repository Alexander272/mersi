import Box from '@mui/material/Box'

import StyledGridOverlay from '../StyledGridOverlay'

import EmptyIcon from './EmptyIcon'

export const NoRowsOverlay = () => {
	return (
		<StyledGridOverlay>
			<EmptyIcon />
			<Box mt={1} color='text.secondary'>
				Не найдено ни одной позиции
			</Box>
		</StyledGridOverlay>
	)
}
