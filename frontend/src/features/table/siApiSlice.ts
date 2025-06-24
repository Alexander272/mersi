import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IGetSiDTO, ISI, ISiDTO } from './types/si'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { buildSiUrlParams } from './utils/buildUrlParams'

const SIApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getSI: builder.query<{ data: ISI[]; total: number }, IGetSiDTO>({
			query: params => ({
				url: API.si.base,
				method: 'GET',
				params: buildSiUrlParams(params),
			}),
			providesTags: [{ type: 'SI', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					console.error(fetchError)
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createSi: builder.mutation<null, ISiDTO>({
			query: body => ({
				url: API.si.base,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'ALL' },
				{ type: 'Instrument', id: 'Unique' },
			],
		}),
	}),
})

export const { useGetSIQuery, useCreateSiMutation } = SIApiSlice
