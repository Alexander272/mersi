import { JSX } from 'react'
import { Stack } from '@mui/material'
import type { DialogProps as MuiDialogProps } from '@mui/material/Dialog'
import MuiDialog from '@mui/material/Dialog'
import MuiDialogActions from '@mui/material/DialogActions'
import MuiDialogContent from '@mui/material/DialogContent'
import MuiDialogContentText from '@mui/material/DialogContentText'
import MuiDialogTitle from '@mui/material/DialogTitle'

interface IDialogProps extends MuiDialogProps {
	title: string
	headerActions?: JSX.Element
	body: string | JSX.Element
	actions?: JSX.Element
}

export const Dialog = (props: IDialogProps) => {
	const { actions, title, body, headerActions, ...other } = props

	return (
		<MuiDialog {...other}>
			<Stack direction={'row'} width={'100%'} justifyContent={'space-between'} alignItems={'center'}>
				<MuiDialogTitle>{title}</MuiDialogTitle>

				{headerActions}
			</Stack>

			<MuiDialogContent>
				{typeof body === 'string' ? <MuiDialogContentText>{body}</MuiDialogContentText> : body}
			</MuiDialogContent>

			{actions && <MuiDialogActions>{actions}</MuiDialogActions>}
		</MuiDialog>
	)
}
