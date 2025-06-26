import { SvgIcon, SxProps, Theme } from '@mui/material'
import { FC } from 'react'

export const LongRightIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg
				xmlns='http://www.w3.org/2000/svg'
				shapeRendering='geometricPrecision'
				textRendering='geometricPrecision'
				imageRendering='optimizeQuality'
				fillRule='evenodd'
				clipRule='evenodd'
				viewBox='0 0 512 243.58'
			>
				<path
					fillRule='nonzero'
					d='M373.57 0 512 120.75 371.53 243.58l-20.92-23.91 94.93-83L0 137.09v-31.75l445.55-.41-92.89-81.02z'
				/>
			</svg>
		</SvgIcon>
	)
}
