import { PayloadAction, createSlice } from '@reduxjs/toolkit'

import type { IDepartment, IEmployee } from './types/employee'
import { RootState } from '@/app/store'

interface IEmployeesState {
	employee?: IEmployee
	department?: IDepartment
}

const initialState: IEmployeesState = {}

const employeesSlice = createSlice({
	name: 'employees',
	initialState,
	reducers: {
		setEmployee: (state, action: PayloadAction<IEmployee | undefined>) => {
			state.employee = action.payload
		},

		setDepartment: (state, action: PayloadAction<IDepartment | undefined>) => {
			state.department = action.payload
		},

		resetEmployeeState: () => initialState,
	},
})

export const getEmployee = (state: RootState) => state.employees.employee
export const getDepartment = (state: RootState) => state.employees.department

export const employeesPath = employeesSlice.name
export const employeesReducer = employeesSlice.reducer

export const { setEmployee, setDepartment, resetEmployeeState } = employeesSlice.actions
