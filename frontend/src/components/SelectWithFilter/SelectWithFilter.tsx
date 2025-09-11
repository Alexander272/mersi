import { FC, useRef, useState } from 'react'
import {
	Autocomplete,
	autocompleteClasses,
	AutocompleteCloseReason,
	Box,
	ClickAwayListener,
	InputAdornment,
	Popper,
	Stack,
	styled,
	SxProps,
	TextField,
	Theme,
	useTheme,
} from '@mui/material'

import { SearchIcon } from '../Icons/SearchIcon'
import { ListboxComponent } from './Listbox'

export type Option = {
	id: string
	name: string
}

type Props = {
	label?: string
	headerLabel?: string
	values: Option[]
	options: Option[]
	disabled?: boolean
	onChange: (values: Option[]) => void
	sx?: SxProps<Theme>
}

export const SelectWithFilter: FC<Props> = ({ label, headerLabel, values, options, onChange, disabled, sx }) => {
	const theme = useTheme()
	const anchor = useRef<HTMLParagraphElement>(null)
	const [open, setOpen] = useState(false)

	const openHandler = () => {
		setOpen(true)
	}
	const closeHandler = () => {
		setOpen(false)
	}

	return (
		<Stack width={'100%'} sx={sx}>
			<TextField
				ref={anchor}
				label={label || 'Значение'}
				value={values.length ? values.map(v => v.name).join(', ') : ''}
				onClick={openHandler}
				disabled={disabled}
				slotProps={{
					htmlInput: {
						readOnly: true,
						sx: { cursor: 'pointer', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' },
					},
					inputLabel: {
						shrink: open || Boolean(values.length),
					},
				}}
			/>

			<Popper
				open={open}
				anchorEl={anchor.current}
				placement='bottom-start'
				sx={{
					width: anchor.current ? anchor.current.clientWidth + 2 : 'auto',
					border: '1px solid #e1e4e8',
					boxShadow: `0 8px 24px rgba(149, 157, 165, 0.2)`,
					color: '#24292e',
					backgroundColor: '#fff',
					borderRadius: 2,
					zIndex: theme.zIndex.modal + 10,
					fontSize: 13,
				}}
			>
				<ClickAwayListener onClickAway={closeHandler}>
					<div>
						{headerLabel && (
							<Box
								sx={t => ({
									borderBottom: '1px solid #30363d',
									padding: '8px 10px',
									fontWeight: 600,
									...t.applyStyles('light', {
										borderBottom: '1px solid #eaecef',
									}),
								})}
							>
								{headerLabel}
							</Box>
						)}
						<Autocomplete
							open
							multiple
							onClose={(_event: React.ChangeEvent<object>, reason: AutocompleteCloseReason) => {
								if (reason === 'escape') {
									closeHandler()
								}
							}}
							value={values}
							onChange={(event, newValue, reason) => {
								if (
									event.type === 'keydown' &&
									((event as React.KeyboardEvent).key === 'Backspace' ||
										(event as React.KeyboardEvent).key === 'Delete') &&
									reason === 'removeOption'
								) {
									return
								}
								onChange(newValue)
							}}
							disableCloseOnSelect
							renderValue={() => null}
							noOptionsText='Ничего не найдено'
							renderOption={(props, option, state) => [props, option, state] as React.ReactNode}
							options={options}
							getOptionLabel={option => option.name}
							renderInput={params => (
								<TextField
									ref={params.InputProps.ref}
									autoFocus
									fullWidth
									placeholder='Поиск'
									sx={{ padding: '8px 16px', borderBottom: '1px solid #eaecef' }}
									slotProps={{
										htmlInput: {
											...params.inputProps,
										},
										input: {
											startAdornment: (
												<InputAdornment position='start'>
													<SearchIcon fontSize={16} ml={1} />
												</InputAdornment>
											),
										},
									}}
								/>
							)}
							slotProps={{ listbox: { component: ListboxComponent } }}
							slots={{
								popper: PopperComponent,
							}}
						/>
					</div>
				</ClickAwayListener>
			</Popper>
		</Stack>
	)
}

interface PopperComponentProps {
	anchorEl?: unknown
	disablePortal?: boolean
	open: boolean
}

function PopperComponent(props: PopperComponentProps) {
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	const { disablePortal, anchorEl, open, ...other } = props
	return <StyledAutocompletePopper {...other} />
}

const StyledAutocompletePopper = styled('div')(({ theme }) => ({
	[`& .${autocompleteClasses.paper}`]: {
		boxShadow: 'none',
		color: 'inherit',
		fontSize: 13,
	},
	[`& .${autocompleteClasses.listbox}`]: {
		padding: 0,
		backgroundColor: '#fff',

		[`& .${autocompleteClasses.option}`]: {
			// flexDirection: 'column',
			// alignItems: 'flex-start',
			padding: 8,
			paddingBottom: 0,
			margin: 6,
			borderRadius: 8,

			[`&.${autocompleteClasses.focused}, &.${autocompleteClasses.focused}[aria-selected="true"]`]: {
				backgroundColor: theme.palette.action.hover,
			},
		},
		'& ul': {
			padding: 0,
			margin: 0,
		},
	},
	[`&.${autocompleteClasses.popperDisablePortal}`]: {
		position: 'relative',
	},
}))
