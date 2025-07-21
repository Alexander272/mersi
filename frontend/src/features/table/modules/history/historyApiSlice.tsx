import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IHistoryType } from './types/historyTypes'
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
	}),
})

export const { useGetHistoryTypesQuery } = historyApiSlice
