import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IRepair, IRepairDTO } from './types/repair'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const repairApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getRepair: builder.query<{ data: IRepair[] }, string>({
			query: instrument => ({
				url: API.si.repair,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'Repair' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createRepair: builder.mutation<null, IRepairDTO>({
			query: body => ({
				url: API.si.repair,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'Repair' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const { useGetRepairQuery, useCreateRepairMutation } = repairApiSlice
