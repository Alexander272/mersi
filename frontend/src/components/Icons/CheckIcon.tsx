import { SvgIcon, SxProps, Theme } from '@mui/material'
import { FC } from 'react'

export const CheckIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 512 457.57'>
				<path d='M0 220.57c100.43-1.33 121-5.2 191.79 81.5 54.29-90 114.62-167.9 179.92-235.86C436-.72 436.5-.89 512 .24 383.54 143 278.71 295.74 194.87 457.57 150 361.45 87.33 280.53 0 220.57z' />
			</svg>
		</SvgIcon>
	)
}
