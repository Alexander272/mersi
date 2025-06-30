import { FC, JSX, MouseEvent } from 'react'
import { CircularProgress, IconButton, ListItemIcon, ListItemText, MenuItem, Tooltip, useTheme } from '@mui/material'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { DialogVariants } from '@/features/dialog/dialogSlice'
import type { IToolsMenu } from '../types/toolsMenu'
import { useToggleFavoriteMutation } from '../toolsMenuApiSlice'
import { VerifyIcon } from '@/components/Icons/VerifyIcon'
import { RepairIcon } from '@/components/Icons/RepairIcon'
import { ToolboxIcon } from '@/components/Icons/ToolboxIcon'
import { ProductReplace } from '@/components/Icons/ProductReplace'
import { ProductReturn } from '@/components/Icons/ProductReturn'
import { DocumentCheckIcon } from '@/components/Icons/DocumentCheckIcon'
import { FileDownloadIcon } from '@/components/Icons/FileDownloadIcon'
import { StarIcon } from '@/components/Icons/StarIcon'

type Props = { onClick?: () => void; item: IToolsMenu; label?: string }
type MenuItem = { el: FC<Props>; action?: DialogVariants }

export const Icons = new Map<string, JSX.Element>([
	['graphic', <DocumentCheckIcon fontSize={20} fill={'#757575'} />],
	['export', <FileDownloadIcon fontSize={20} fill={'#757575'} />],
	['verification', <VerifyIcon fontSize={18} fill={'#757575'} />],
	['repair-info', <RepairIcon fontSize={18} fill={'#757575'} />],
	['preservation-info', <ToolboxIcon fontSize={18} fill={'#757575'} />],
	['transfer-to-save', <ProductReplace fontSize={18} fill={'#757575'} />],
	['transfer-to-department', <ProductReturn fontSize={18} fill={'#757575'} />],
])

export const ToolsItem: FC<Props> = ({ onClick, item, label }) => {
	const { palette } = useTheme()

	const [toggle, { isLoading }] = useToggleFavoriteMutation()

	const favoriteHandler = async (event: MouseEvent<HTMLButtonElement>) => {
		event.stopPropagation()

		try {
			await toggle({ id: item.id, favorite: !item.favorite }).unwrap()
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	}

	return (
		<MenuItem onClick={onClick} sx={{ height: 38 }}>
			<ListItemIcon>{Icons.get(item.name)}</ListItemIcon>
			<ListItemText>{label ? label : item.label}</ListItemText>

			{item.canBeFavorite && (
				<Tooltip title='Добавить в контекстное меню' enterDelay={600}>
					<IconButton onClick={favoriteHandler} sx={{ marginY: -0.5 }}>
						{isLoading ? (
							<CircularProgress size={18} />
						) : (
							<StarIcon
								fontSize={18}
								fill={item.favorite ? palette.primary.main : palette.action.disabled}
							/>
						)}
					</IconButton>
				</Tooltip>
			)}
		</MenuItem>
	)
}

export const MenuItems = new Map<string, MenuItem>([
	['export', { el: ToolsItem }],
	['graphic', { el: ToolsItem }],
	['verification', { el: ToolsItem }],
	['repair-info', { el: ToolsItem, action: 'AddRepair' }],
	['preservation-info', { el: ToolsItem, action: 'AddPreservation' }],
	['transfer-to-save', { el: ToolsItem, action: 'AddTransferToSave' }],
	['transfer-to-department', { el: ToolsItem, action: 'AddTransferToDep' }],
])
