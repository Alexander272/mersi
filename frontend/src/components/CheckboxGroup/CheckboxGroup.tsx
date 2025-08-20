import { ChangeEvent, FC, useState } from 'react'
import { Checkbox, Collapse, FormControlLabel, IconButton, Stack, SxProps, Theme, useTheme } from '@mui/material'

import { LeftArrowIcon } from '../Icons/LeftArrowIcon'

type Item = {
	id: string
	name: string
}

type Props = {
	checked: Map<string, boolean>
	data: { name: string; list: Item[] }
	onChange: (checked: Map<string, boolean>) => void
	// data: Map<string, Item[]>
	sx?: SxProps<Theme>
}

export const CheckboxGroup: FC<Props> = ({ checked, data, onChange, sx }) => {
	const { palette } = useTheme()
	const [open, setOpen] = useState(true)

	const ids = data.list.reduce((a, v) => [...a, v.id], [] as string[])

	const toggleHandler = () => setOpen(prev => !prev)

	const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
		console.log(event.target.name, event.target.checked)
		const tmp = new Map(checked)
		if (event.target.name != data.name) {
			tmp.set(event.target.name, event.target.checked)
			onChange(tmp)
			return
		}

		ids.forEach(id => {
			tmp.set(id, event.target.checked)
		})
		onChange(tmp)
	}

	return (
		<Stack sx={sx} mb={1}>
			<Stack direction={'row'} spacing={2} justifyContent={'space-between'} borderBottom={'1px solid #e0e0e0'}>
				<FormControlLabel
					label={data.name}
					control={
						<Checkbox
							name={data.name}
							checked={ids.every(id => checked.get(id) === true) || false}
							indeterminate={ids.some(id => checked.get(id) != checked.get(ids[0]))}
							onChange={handleChange}
						/>
					}
					sx={{
						pl: 1,
						width: '100%',
						transition: 'all 0.3s ease-in-out',
						borderRadius: 3,
						':hover': { backgroundColor: palette.action.hover },
					}}
				/>

				<IconButton onClick={toggleHandler} sx={{ width: 42 }}>
					<LeftArrowIcon fontSize={14} transform={open ? 'rotate(90deg)' : 'rotate(-90deg)'} />
				</IconButton>
			</Stack>

			<Collapse in={open}>
				<Stack spacing={1} pt={1} sx={{ background: '#f5f5f57d', borderRadius: '0 0 8px 8px' }}>
					{data.list.map(d => (
						<FormControlLabel
							key={d.id}
							label={d.name}
							control={
								<Checkbox name={d.id} checked={checked.get(d.id) || false} onChange={handleChange} />
							}
							sx={{
								transition: 'all 0.3s ease-in-out',
								pl: 4,
								borderRadius: 3,
								':hover': { backgroundColor: palette.action.hover },
							}}
						/>
					))}
				</Stack>
			</Collapse>
		</Stack>
	)
}
