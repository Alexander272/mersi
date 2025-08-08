import { FC, useEffect } from 'react'
import { Divider, FormControlLabel, Stack, Switch } from '@mui/material'
import { Controller, FieldArrayWithId, useFieldArray, useFormContext } from 'react-hook-form'

import type { IColumn } from '@/features/sections/modules/columns/types/columns'
import { Item } from './Item'

type Props = {
	index: number
	data: FieldArrayWithId<{ data: IColumn[] }, 'data', 'id'>
}

export const Group: FC<Props> = ({ index, data }) => {
	const path = `data.${index}.children` as const
	const { control, watch, setValue } = useFormContext()
	const { fields } = useFieldArray({ control, name: path })
	const hidden = watch(`data.${index}.hidden`)

	useEffect(() => {
		if (hidden == undefined) return
		data.children?.forEach((_c, i) => {
			setValue(`data.${index}.children.${i}.hidden`, hidden)
		})
	}, [data.children, hidden, index, setValue])

	return (
		<Stack mb={1} px={1} pt={1} sx={{ border: '1px solid #e0e0e0', borderRadius: 2 }}>
			<Controller
				control={control}
				name={`data.${index}.hidden`}
				render={({ field }) => (
					<FormControlLabel
						label={data.name}
						sx={{
							color: !field.value ? 'inherit' : '#505050',
							transition: '.2s color ease-in-out',
							userSelect: 'none',
							mb: 0.5,
						}}
						control={
							<Switch checked={!field.value} onChange={event => field.onChange(!event.target.checked)} />
						}
					/>
				)}
			/>
			<Divider />

			<Stack pt={1}>
				{fields.map((f, i) => (
					<Item
						key={f.id}
						index={i}
						label={(f as FieldArrayWithId<{ data: IColumn[] }, 'data', 'id'>).name}
						path={path}
					/>
				))}
			</Stack>
		</Stack>
	)
}
