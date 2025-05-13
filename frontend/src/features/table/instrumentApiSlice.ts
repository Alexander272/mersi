import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const instrumentApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getUniqueInstrumentData: builder.query<{ data: string[] }, string>({
			query: field => ({
				url: `${API.si.instruments.unique}/${field}`,
			}),
			providesTags: [{ type: 'Instrument', id: 'Unique' }],
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

export const { useGetUniqueInstrumentDataQuery } = instrumentApiSlice
