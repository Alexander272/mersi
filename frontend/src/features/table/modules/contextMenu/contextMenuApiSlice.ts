import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { IContextMenu, IContextMenuDTO } from './types/context'

const contextMenuApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getContextMenu: builder.query<{ data: IContextMenu[] }, { section: string; isFull?: boolean }>({
			query: req => ({
				url: API.si.context,
				// params: new URLSearchParams({ section: req.section, isFull: req.isFull ? 'true' : 'false' }),
				params: new URLSearchParams([['section', req.section], req.isFull ? ['isFull', 'true'] : []]),
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

		createContextMenu: builder.mutation<{ message: string }, IContextMenuDTO>({
			query: data => ({
				url: API.si.context,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Context' }],
		}),
		updateContextMenu: builder.mutation<{ message: string }, IContextMenuDTO>({
			query: data => ({
				url: `${API.si.context}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Context' }],
		}),
		deleteContextMenu: builder.mutation<{ message: string }, string>({
			query: id => ({
				url: `${API.si.context}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Context' }],
		}),
	}),
})

export const {
	useGetContextMenuQuery,
	useCreateContextMenuMutation,
	useUpdateContextMenuMutation,
	useDeleteContextMenuMutation,
} = contextMenuApiSlice
