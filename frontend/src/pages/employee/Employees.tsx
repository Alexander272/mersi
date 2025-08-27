import { useState } from 'react'
import { Box, Breadcrumbs, Stack } from '@mui/material'

import { DepartmentList } from '@/features/departments/components/DepartmentList'
import { EmployeeTable } from '@/features/employees/components/EmployeeTable'
import { DepartmentForm } from '@/features/departments/components/DepartmentForm'
import { PageBox } from '@/components/PageBox/PageBox'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { AppRoutes } from '../router/routes'

// страница для управления работниками и департаментами
export default function Employees() {
	const [department, setDepartment] = useState('new')

	const departmentHandler = (department: string) => {
		setDepartment(department)
	}

	return (
		<PageBox>
			<Box
				borderRadius={3}
				padding={2}
				margin={'0 auto'}
				width={{ xl: '66%', lg: '86%', md: '100%' }}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				flexGrow={1}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: 2 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.Employees} active>
						Подразделения
					</Breadcrumb>
				</Breadcrumbs>

				<Stack direction={'row'} spacing={2} height={'100%'}>
					<DepartmentList department={department} setDepartment={departmentHandler} />
					<Stack width={'100%'} spacing={3} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
						<DepartmentForm department={department} setDepartment={departmentHandler} />
						{department != 'new' && <EmployeeTable department={department} hasResponsible />}
					</Stack>
				</Stack>
			</Box>
		</PageBox>
	)
}
