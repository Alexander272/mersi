import { toast } from 'react-toastify'

import type { IEmployee } from './types/employee'
import type { IBaseFetchError } from '@/app/types/error'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const employeesApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getEmployees: builder.query<{ data: IEmployee[] }, string | null>({
			query: departmentId => ({
				url: `${API.employees}`,
				params: new URLSearchParams(departmentId ? { departmentId } : undefined),
			}),
			providesTags: (_result, _error, arg) => [{ type: 'Employees', id: arg || 'all' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getUniqueEmployee: builder.query<{ data: IEmployee[] }, null>({
			query: () => `${API.employees}/unique`,
			providesTags: [{ type: 'Employees', id: 'unique' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getEmployeeById: builder.query<{ data: IEmployee }, string>({
			query: id => `${API.employees}/${id}`,
			providesTags: (_result, _error, arg) => [
				{ type: 'Employees', id: arg },
				{ type: 'Employees', id: 'all' },
			],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createEmployee: builder.mutation<null, IEmployee>({
			query: data => ({
				url: `${API.employees}`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: (_result, _error, arg) => [
				{ type: 'Employees', id: 'all' },
				{ type: 'Employees', id: 'unique' },
				{ type: 'Employees', id: arg.departmentId },
			],
		}),

		updateEmployee: builder.mutation<null, IEmployee>({
			query: data => ({
				url: `${API.employees}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: (_result, _error, arg) => [
				{ type: 'Employees', id: 'all' },
				{ type: 'Employees', id: 'unique' },
				{ type: 'Employees', id: arg.departmentId },
			],
		}),

		deleteEmployee: builder.mutation<null, IEmployee>({
			query: data => ({
				url: `${API.employees}/${data.id}`,
				method: 'DELETE',
			}),
			invalidatesTags: (_result, _error, arg) => [
				{ type: 'Employees', id: 'all' },
				{ type: 'Employees', id: 'unique' },
				{ type: 'Employees', id: arg.departmentId },
			],
		}),
	}),
})

export const {
	useGetEmployeesQuery,
	useGetUniqueEmployeeQuery,
	useGetEmployeeByIdQuery,
	useCreateEmployeeMutation,
	useUpdateEmployeeMutation,
	useDeleteEmployeeMutation,
} = employeesApiSlice
