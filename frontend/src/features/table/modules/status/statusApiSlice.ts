import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IStatus, IStatusDTO } from './types/status'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const statusApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getStatuses: builder.query<{ data: IStatus[] }, string>({
			query: section => ({
				url: API.status,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Status', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					console.log(error)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createStatuses: builder.mutation<{ message: string }, IStatusDTO[]>({
			query: data => ({
				url: `${API.status}/several`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Status', id: 'ALL' }],
		}),
		updateStatuses: builder.mutation<{ message: string }, IStatusDTO[]>({
			query: data => ({
				url: `${API.status}/several`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'Status', id: 'ALL' }],
		}),
		deleteStatuses: builder.mutation<{ message: string }, string[]>({
			query: data => ({
				url: `${API.status}/several`,
				method: 'DELETE',
				body: data,
			}),
			invalidatesTags: [{ type: 'Status', id: 'ALL' }],
		}),
	}),
})

export const { useGetStatusesQuery, useCreateStatusesMutation, useUpdateStatusesMutation, useDeleteStatusesMutation } =
	statusApiSlice
