import { FC } from 'react'
import { SvgIcon, SxProps, Theme } from '@mui/material'

export const DocIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg viewBox='0 0 102.55 122.88' xmlSpace='preserve'>
				<path d='M102.55 122.88H0V0h77.66l24.89 26.43v96.45zM64.76 58.39h4.66l-5.38 23.86-6.99-33.55H45.81l-8.14 33.55-6.1-33.55H19.82l11.93 54h11.49l8.05-31.94 7.2 31.94h10.59l13.67-54H64.76v9.69zm31.37 57.59V30.44H73.44V5.91H6.51v110.07h89.62z' />
			</svg>
		</SvgIcon>
	)
}
