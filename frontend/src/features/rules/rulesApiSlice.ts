import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { IRuleItem } from './types/rules'

const rulesApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getRules: builder.query<{ data: IRuleItem[] }, null>({
			query: () => ({
				url: API.rules,
			}),
			providesTags: [{ type: 'SI', id: 'Rules' }],
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

export const { useGetRulesQuery } = rulesApiSlice
