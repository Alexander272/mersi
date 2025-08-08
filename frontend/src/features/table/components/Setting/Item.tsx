import { FC } from 'react'
import { FormControlLabel, InputAdornment, Stack, Switch, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'
import { WidthIcon } from '@/components/Icons/WidthIcon'

type Props = {
	index: number
	path?: string
	label: string
}

export const Item: FC<Props> = ({ index, path = 'data', label }) => {
	const { control } = useFormContext()

	return (
		<Stack direction={'row'} alignItems={'center'} mb={1}>
			<Controller
				control={control}
				name={`${path}.${index}.hidden`}
				render={({ field }) => (
					<FormControlLabel
						label={label}
						sx={{
							color: !field.value ? 'inherit' : '#505050',
							transition: '.2s color ease-in-out',
							userSelect: 'none',
						}}
						control={
							<Switch checked={!field.value} onChange={event => field.onChange(!event.target.checked)} />
						}
					/>
				)}
			/>
			<Controller
				control={control}
				name={`${path}.${index}.width`}
				rules={{
					min: 50,
					max: 800,
					pattern: {
						value: /^(0|[1-9]\d*)(\.\d+)?$/,
						message: '',
					},
				}}
				render={({ field, fieldState: { error } }) => (
					<TextField
						{...field}
						error={Boolean(error)}
						helperText={error ? '50 ≤ Размер ≤ 800' : ''}
						slotProps={{
							input: {
								type: 'number',
								startAdornment: (
									<InputAdornment position='start'>
										<WidthIcon fontSize={18} />
									</InputAdornment>
								),
							},
						}}
						sx={{
							width: 120,
							minWidth: 120,
							ml: 'auto',
							'.MuiFormHelperText-root': { mx: 0.8 },
						}}
					/>
				)}
			/>
		</Stack>
	)
}
