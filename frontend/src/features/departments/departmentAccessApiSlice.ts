import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IChangeDepAccess, IDepartmentAccess } from './types/departments'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

export const departmentAccessApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getDepAccesses: builder.query<{ data: IDepartmentAccess[] }, { department?: string }>({
			query: req => ({
				url: `${API.departmentAccesses}/${req.department}`,
				method: 'GET',
			}),
			providesTags: [{ type: 'DepartmentAccesses', id: 'All' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		// getResponsibleByUser: builder.query<{ data: IResponsible[] }, null>({
		//     query: () => `${API.responsible}/sso`,
		//     providesTags: [{ type: 'Responsible', id: 'all' }],
		//     onQueryStarted: async (_arg, api) => {
		//         try {
		//             await api.queryFulfilled
		//         } catch (error) {
		//             const fetchError = (error as IBaseFetchError).error
		//             toast.error(fetchError.data.message, { autoClose: false })
		//         }
		//     },
		// }),

		changeDepAccesses: builder.mutation<null, IChangeDepAccess>({
			query: data => ({
				url: `${API.departmentAccesses}/replace`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'DepartmentAccesses', id: 'All' }],
		}),
	}),
})

export const { useGetDepAccessesQuery, useChangeDepAccessesMutation } = departmentAccessApiSlice
