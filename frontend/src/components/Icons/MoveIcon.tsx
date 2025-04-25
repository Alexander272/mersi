import { FC } from 'react'
import { SvgIcon, SxProps, Theme } from '@mui/material'

export const MoveIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 122.88 80.59' xmlSpace='preserve'>
				<path d='M122.88 27.53 94.25 0v17.72h-22.7v19.62h22.7v17.72l28.63-27.53zM0 53.06l28.63 27.53V62.87h22.7V43.25h-22.7V25.53L0 53.06zm54.84 9.81H68.4V43.25H54.84v19.62zm16.8 0h7.7V43.25h-7.7v19.62zm-20.4-45.15h-7.7v19.62h7.7V17.72zm16.8 0H54.49v19.62h13.56V17.72h-.01z' />
			</svg>
		</SvgIcon>
	)
}
