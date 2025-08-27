import { FC, useState } from 'react'
import { Autocomplete, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { ICreateFormField } from '@/features/sections/modules/form/types/create'
import { useAppSelector } from '@/hooks/redux'
import { useLazyGetUniqueInstrumentDataQuery } from '../../instrumentApiSlice'
import { getSection } from '@/features/sections/sectionSlice'

type Props = {
	data: ICreateFormField
}

export const AutocompleteField: FC<Props> = ({ data }) => {
	const [options, setOptions] = useState<string[]>([])
	const section = useAppSelector(getSection)

	const { control, watch } = useFormContext()
	const watchField = watch(data.hide)

	// const { data: options, isFetching } = useGetUniqueInstrumentDataQuery(data.field, { skip: !data.field })

	const [getUnique, { isLoading }] = useLazyGetUniqueInstrumentDataQuery()

	const focusHandler = async () => {
		if (!section) return
		const res = await getUnique({ field: data.field, section: section?.id || '' }).unwrap()
		setOptions(res.data || [])
	}

	if (data.hide && watchField) return null
	return (
		<Controller
			name={data.path + '.' + data.field}
			control={control}
			rules={{ required: data.isRequired }}
			render={({ field: { onChange, value, ref }, fieldState: { error } }) => (
				<Autocomplete
					value={value || ''}
					freeSolo
					disableClearable
					autoComplete
					options={options}
					loading={isLoading}
					loadingText='Поиск похожих значений...'
					noOptionsText='Ничего не найдено'
					onChange={(_event, value) => {
						onChange(value)
					}}
					onFocus={focusHandler}
					renderInput={params => (
						<TextField
							{...params}
							label={data.fieldName}
							onChange={onChange}
							error={Boolean(error)}
							helperText={error?.message}
							inputRef={ref}
						/>
					)}
				/>
			)}
		/>
	)
}
