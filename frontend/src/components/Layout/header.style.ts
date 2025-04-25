import { styled } from '@mui/material'
import { Link } from 'react-router'

export const NavLink = styled(Link)(({ theme }) => ({
	textDecoration: 'none',
	fontSize: '1.2rem',
	color: theme.palette.primary.main,
	height: '100%',
	display: 'flex',
	alignItems: 'center',
	justifyContent: 'center',
	position: 'relative',

	'&:after': {
		content: '""',
		position: 'absolute',
		bottom: 0,
		height: 2,
		width: '100%',
		maxWidth: '0%',
		backgroundColor: theme.palette.primary.main,
		transition: 'all 0.3s ease-in-out',
	},

	'&:hover': {
		'&:after': {
			maxWidth: '100%',
		},
	},
}))

export const NavButton = styled('div')(({ theme }) => ({
	cursor: 'pointer',
	fontSize: '1.2rem',
	color: theme.palette.primary.main,
	height: '100%',
	display: 'flex',
	alignItems: 'center',
	justifyContent: 'center',
	position: 'relative',

	'&:after': {
		content: '""',
		position: 'absolute',
		bottom: 0,
		height: 2,
		width: '100%',
		maxWidth: '0%',
		backgroundColor: theme.palette.primary.main,
		transition: 'all 0.3s ease-in-out',
	},

	'&:hover': {
		'&:after': {
			maxWidth: '100%',
		},
	},
}))

// export const NavButton = styled.p`
// 	cursor: pointer;
// 	font-size: 1.2rem;
// 	color: var(--primary-color);
// 	height: 100%;
// 	display: flex;
// 	align-items: center;
// 	justify-content: center;
// 	position: relative;

// 	&:after {
// 		content: '';
// 		position: absolute;
// 		bottom: 0px;
// 		height: 2px;
// 		width: 100%;
// 		max-width: 0%;
// 		background-color: var(--primary-color);
// 		transition: all 0.3s ease-in-out;
// 	}

// 	&:hover {
// 		&:after {
// 			max-width: 100%;
// 		}
// 	}
// `
