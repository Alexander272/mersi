import { FC } from 'react'
import { SvgIcon, SxProps, Theme } from '@mui/material'

export const ChangeArrowsIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg
				xmlns='http://www.w3.org/2000/svg'
				shapeRendering='geometricPrecision'
				textRendering='geometricPrecision'
				imageRendering='optimizeQuality'
				fillRule='evenodd'
				clipRule='evenodd'
				viewBox='0 0 512 421.69'
			>
				<path
					fillRule='nonzero'
					d='M225.95 233.54v-96.9c0-3.72 3.02-6.74 6.74-6.74H396.1c3.72 0 6.74 3.02 6.74 6.74v40.92l92.05-78.15-92.05-78.15v40.82c0 3.72-3.02 6.74-6.74 6.74H164.97v164.72h60.98zm13.47-90.17v96.91c0 3.72-3.01 6.74-6.73 6.74h-74.46c-3.72 0-6.74-3.02-6.74-6.74V62.08c0-3.72 3.02-6.74 6.74-6.74h231.14V6.72c.01-1.53.53-3.08 1.6-4.34a6.731 6.731 0 0 1 9.48-.79l109.09 92.62c.31.26.6.54.87.86 2.4 2.83 2.04 7.08-.79 9.47l-108.84 92.4a6.732 6.732 0 0 1-11.41-4.84v-48.73H239.42zm46.63 44.78v96.9c0 3.72-3.02 6.74-6.74 6.74H115.9c-3.72 0-6.74-3.02-6.74-6.74v-40.92l-92.05 78.15 92.05 78.15v-40.82c0-3.72 3.02-6.74 6.74-6.74h231.13V188.15h-60.98zm-13.47 90.17v-96.91c0-3.72 3.01-6.73 6.73-6.73h74.46c3.72 0 6.74 3.01 6.74 6.73v178.2c0 3.72-3.02 6.74-6.74 6.74H122.63v48.62a6.757 6.757 0 0 1-1.6 4.34 6.731 6.731 0 0 1-9.48.79L2.46 327.48c-.31-.26-.6-.54-.87-.86-2.4-2.83-2.04-7.08.79-9.47l108.84-92.4a6.732 6.732 0 0 1 11.41 4.84v48.73h149.95z'
				/>
			</svg>
		</SvgIcon>
	)
}
