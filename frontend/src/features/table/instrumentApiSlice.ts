import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { IInstrumentDTO } from './types/instrument'

const instrumentApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getInstrumentById: builder.query<{ data: IInstrumentDTO }, string>({
			query: id => `${API.si.instruments.base}/${id}`,
			providesTags: (_res, _err, arg) => [{ type: 'Instrument', id: arg }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getUniqueInstrumentData: builder.query<{ data: string[] }, { field: string; section: string }>({
			query: data => ({
				url: `${API.si.instruments.unique}/${data.field}`,
				params: new URLSearchParams({ section: data.section }),
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

		changeStatus: builder.mutation<{ data: IInstrumentDTO }, { id: string; status: string }[]>({
			query: data => ({
				url: `${API.si.instruments.base}/status`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'SI', id: 'ALL' }],
			// invalidatesTags: (_res, _err, arg) => [
			// 	{ type: 'Instrument', id: arg.id },
			// 	{ type: 'SI', id: 'ALL' },
			// ],
		}),
	}),
})

export const {
	useGetInstrumentByIdQuery,
	useLazyGetInstrumentByIdQuery,
	useGetUniqueInstrumentDataQuery,
	useLazyGetUniqueInstrumentDataQuery,
	useChangeStatusMutation,
} = instrumentApiSlice
