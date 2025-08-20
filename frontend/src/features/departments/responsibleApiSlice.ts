import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IChangeResponsible, IResponsible } from './types/responsible'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

export const responsibleApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getResponsible: builder.query<{ data: IResponsible[] }, { department?: string; sso?: string }>({
			query: req => ({
				url: API.responsible,
				method: 'GET',
				params: { department: req.department, sso: req.sso },
			}),
			providesTags: [{ type: 'Responsible', id: 'all' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		getResponsibleByUser: builder.query<{ data: IResponsible[] }, null>({
			query: () => `${API.responsible}/sso`,
			providesTags: [{ type: 'Responsible', id: 'all' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		changeResponsible: builder.mutation<null, IChangeResponsible>({
			query: data => ({
				url: `${API.responsible}/change`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Responsible', id: 'all' }],
		}),
	}),
})

export const { useGetResponsibleQuery, useGetResponsibleByUserQuery, useChangeResponsibleMutation } =
	responsibleApiSlice
