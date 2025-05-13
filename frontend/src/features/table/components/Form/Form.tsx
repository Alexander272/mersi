import { FC } from 'react'
import { Stack } from '@mui/material'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { TextField } from './TextField'
import { NumberField } from './NumberField'
import { DateField } from './DateField'
import { FileField } from './FileField'
import { AutocompleteField } from './AutocompleteField'
import { SelectField } from './SelectField'

type Props = {
	data: ICreateFormField[]
}

export const Form: FC<Props> = ({ data }) => {
	const renderFields = () => {
		return data.map(item => {
			switch (item.type) {
				case 'text':
					return <TextField key={item.id} data={item} />
				case 'number':
					return <NumberField key={item.id} data={item} />
				case 'date':
					return <DateField key={item.id} data={item} />
				case 'file':
					return <FileField key={item.id} data={item} />
				case 'list':
					return <SelectField key={item.id} data={item} />
				case 'autocomplete':
					return <AutocompleteField key={item.id} data={item} />
				// TODO надо еще придумать как выводить поля если они зависят друг от друга
				default:
					return null
			}
		})
	}

	return (
		<Stack spacing={1.5} mb={2}>
			{renderFields()}
		</Stack>
	)
}
