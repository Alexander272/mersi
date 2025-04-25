import { SvgIcon, SxProps, Theme } from '@mui/material'
import { FC } from 'react'

export const ImageIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 102.55 122.88' xmlSpace='preserve'>
				<path d='M102.55 122.88H0V0h77.66l24.89 26.44v96.44zM29.82 45.51c4.33 0 7.84 3.5 7.84 7.83a7.836 7.836 0 0 1-15.67 0 7.83 7.83 0 0 1 7.83-7.83zm32.89 29.76 14.02-24.24 2.07-3.57 1.52 3.84 9.45 44.84h-75.3v-9.12l2.87-.1 5.26-.27 6.06-14.81 3.77.2 2.84 9.98h6.92l8-20.59 1.43-3.69 2.12 3.33 8.97 14.2zm33.42 40.71V32.36H73.45V5.91H6.51v110.07h89.62z' />
			</svg>
		</SvgIcon>
	)
}
