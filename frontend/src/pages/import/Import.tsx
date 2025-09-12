import { Box, Breadcrumbs } from '@mui/material'

import { PermRules } from '@/constants/permissions'
import { useCheckPermission } from '@/features/user/hooks/check'
import { Form } from '@/features/import/components/Form'
import { PageBox } from '@/components/PageBox/PageBox'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { Forbidden } from '../forbidden/ForbiddenLazy'

export default function Import() {
	if (!useCheckPermission(PermRules.Import.Write)) return <Forbidden />
	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				width={'100%'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				height={'fit-content'}
				minHeight={600}
				maxHeight={800}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs sx={{ mb: -3, zIndex: 50, width: 'fit-content' }}>
					<Breadcrumb to={'/'}>Главная</Breadcrumb>
					<Breadcrumb to={'/import'} active>
						Импорт
					</Breadcrumb>
				</Breadcrumbs>

				<Form />
			</Box>
		</PageBox>
	)
}
