import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { IContextMenu } from './types/context'

const contextMenuApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getContextMenu: builder.query<{ data: IContextMenu[] }, string>({
			query: section => ({
				url: API.si.context,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Sections', id: 'Context' }],
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

export const { useGetContextMenuQuery } = contextMenuApiSlice
