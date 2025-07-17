import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IWriteOff, IWriteOffDTO } from './types/writeoff'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const writeOffApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getWriteOff: builder.query<{ data: IWriteOff[] }, string>({
			query: instrument => ({
				url: API.si.writeOff,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'WriteOff' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createWriteOff: builder.mutation<IWriteOff, IWriteOffDTO>({
			query: data => ({
				url: API.si.writeOff,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'WriteOff' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const { useGetWriteOffQuery, useCreateWriteOffMutation } = writeOffApiSlice
