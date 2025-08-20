import { ButtonBase, ButtonBaseProps, styled } from '@mui/material'

const Button = styled(ButtonBase)<ButtonBaseProps>(({ theme }) => ({
	flexGrow: 1,
	width: '100%',
	padding: '6px 20px',
	border: `1px solid ${theme.palette.action.active}`,
	borderRadius: 12,
	transition: 'all .3s ease-in-out',
	':hover': {
		backgroundColor: '#0000000a',
	},
}))

export default Button
