import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IHistoryType, IHistoryTypeDTO } from './types/historyTypes'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const historyApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getHistoryTypes: builder.query<{ data: IHistoryType[] }, string>({
			query: section => ({
				url: API.si.historyTypes,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'SI', id: 'History' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createHistoryTypes: builder.mutation<{ message: string }, IHistoryTypeDTO[]>({
			query: data => ({
				url: `${API.si.historyTypes}/several`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'SI', id: 'History' }],
		}),
		updateHistoryTypes: builder.mutation<{ message: string }, IHistoryTypeDTO[]>({
			query: data => ({
				url: `${API.si.historyTypes}/several`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'SI', id: 'History' }],
		}),
		deleteHistoryTypes: builder.mutation<{ message: string }, string[]>({
			query: data => ({
				url: `${API.si.historyTypes}/several`,
				method: 'DELETE',
				body: data, //TODO это вроде как плохая практика
			}),
			invalidatesTags: [{ type: 'SI', id: 'History' }],
		}),
	}),
})

export const {
	useGetHistoryTypesQuery,
	useCreateHistoryTypesMutation,
	useUpdateHistoryTypesMutation,
	useDeleteHistoryTypesMutation,
} = historyApiSlice
