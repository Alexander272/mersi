import type { FC } from 'react'
import { SvgIcon, type SxProps, type Theme } from '@mui/material'

export const ImageIcon: FC<SxProps<Theme>> = style => {
	return (
		<SvgIcon sx={style}>
			<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 256 256' xmlSpace='preserve'>
				<g
					stroke='none'
					strokeWidth={0}
					strokeDasharray={'none'}
					strokeLinecap='butt'
					strokeLinejoin='miter'
					strokeMiterlimit={10}
					fill='none'
					fillRule='nonzero'
					opacity={1}
				>
					<path
						d='M79.285 74.165h-68.57a4.896 4.896 0 0 1-4.891-4.891V20.725a4.896 4.896 0 0 1 4.891-4.89h68.569a4.897 4.897 0 0 1 4.891 4.89v48.549a4.897 4.897 0 0 1-4.89 4.891zm-68.57-56.33a2.894 2.894 0 0 0-2.891 2.89v48.549a2.894 2.894 0 0 0 2.891 2.891h68.569a2.894 2.894 0 0 0 2.891-2.891V20.725a2.894 2.894 0 0 0-2.891-2.89H10.715z'
						stroke='none'
						strokeWidth={1}
						strokeDasharray={'none'}
						strokeLinecap='butt'
						strokeLinejoin='miter'
						strokeMiterlimit={10}
						fill={(style as { fill: string })?.fill || '#000'}
						// fill='#000'
						fillRule='nonzero'
						opacity={1}
						transform='matrix(2.81 0 0 2.81 1.407 1.407)'
					/>
					<path
						d='M79.285 74.165h-68.57a4.896 4.896 0 0 1-4.891-4.891c0-.235.083-.463.234-.643l19.3-23.011a.998.998 0 0 1 .733-.357.974.974 0 0 1 .755.308l10.698 11.16 16.508-20.226a1 1 0 0 1 .755-.368.93.93 0 0 1 .77.338l28.349 32.137c.161.183.25.418.25.661a4.899 4.899 0 0 1-4.891 4.892zm-71.44-4.551a2.895 2.895 0 0 0 2.871 2.551h68.569a2.894 2.894 0 0 0 2.869-2.537L54.855 38.683 38.379 58.871a.999.999 0 0 1-.735.367.974.974 0 0 1-.762-.307L26.174 47.76 7.845 69.614zM37.315 37.9c-3.577 0-6.487-2.911-6.487-6.488s2.91-6.487 6.487-6.487 6.487 2.91 6.487 6.487-2.91 6.488-6.487 6.488zm0-10.974a4.492 4.492 0 0 0-4.487 4.487 4.492 4.492 0 0 0 4.487 4.488 4.492 4.492 0 0 0 4.487-4.488 4.492 4.492 0 0 0-4.487-4.487zM6.824 52.047a1 1 0 0 1-.951-.691L.241 34.02a4.897 4.897 0 0 1 3.14-6.163l3.135-1.019a1 1 0 0 1 1.309.951v23.257a.999.999 0 0 1-1.001 1.001zm-.999-22.881-1.826.593a2.873 2.873 0 0 0-1.683 1.437 2.866 2.866 0 0 0-.173 2.205l3.682 11.331V29.166zM75.989 17.835h-35.45a1 1 0 0 1-.309-1.951l28.364-9.216a4.895 4.895 0 0 1 6.162 3.14l2.184 6.719a.998.998 0 0 1-.951 1.308zm-29.136-2h27.759l-1.758-5.41a2.891 2.891 0 0 0-3.643-1.855l-22.358 7.265zM83.175 63.211a1 1 0 0 1-1-1V38.954a1 1 0 0 1 1.951-.309l5.633 17.337v-.002a4.856 4.856 0 0 1-.294 3.731 4.859 4.859 0 0 1-2.847 2.431l-3.134 1.02a1.023 1.023 0 0 1-.309.049zm1-17.943v15.566L86 60.24a2.864 2.864 0 0 0 1.683-1.437 2.869 2.869 0 0 0 .174-2.204v-.001l-3.682-11.33z'
						stroke='none'
						strokeWidth={1}
						strokeDasharray={'none'}
						strokeLinecap='butt'
						strokeLinejoin='miter'
						strokeMiterlimit={10}
						fill={(style as { fill: string })?.fill || '#000'}
						fillRule='nonzero'
						opacity={1}
						transform='matrix(2.81 0 0 2.81 1.407 1.407)'
					/>
					<path
						d='M19.899 83.572a4.898 4.898 0 0 1-4.655-3.38l-2.183-6.719a.998.998 0 0 1 .951-1.308h35.449a1 1 0 0 1 .308 1.951l-28.364 9.217a4.873 4.873 0 0 1-1.506.239zm-4.511-9.407 1.758 5.41a2.895 2.895 0 0 0 3.642 1.855l22.36-7.266h-27.76z'
						stroke='none'
						strokeWidth={1}
						strokeDasharray={'none'}
						strokeLinecap='butt'
						strokeLinejoin='miter'
						strokeMiterlimit={10}
						fill={(style as { fill: string })?.fill || '#000'}
						fillRule='nonzero'
						opacity={1}
						transform='matrix(2.81 0 0 2.81 1.407 1.407)'
					/>
				</g>
			</svg>
		</SvgIcon>
	)
}
