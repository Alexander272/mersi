import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IToggleToolsMenu, IToolsMenu, IToolsMenuDTO } from './types/toolsMenu'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const toolsMenuApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getToolsMenu: builder.query<{ data: IToolsMenu[] }, { section: string; isFull?: boolean }>({
			query: req => ({
				url: API.si.tools,
				params: new URLSearchParams([['section', req.section], req.isFull ? ['isFull', 'true'] : []]),
			}),
			providesTags: [{ type: 'Sections', id: 'Tools' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		toggleFavorite: builder.mutation<null, IToggleToolsMenu>({
			query: body => ({
				url: `${API.si.tools}/favorite`,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'Sections', id: 'Tools' },
				{ type: 'Sections', id: 'Context' },
			],
		}),

		createToolsMenu: builder.mutation<null, IToolsMenuDTO>({
			query: data => ({
				url: API.si.tools,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Tools' }],
		}),
		updateToolsMenu: builder.mutation<null, IToolsMenuDTO>({
			query: data => ({
				url: `${API.si.tools}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Tools' }],
		}),
		deleteToolsMenu: builder.mutation<null, string>({
			query: id => ({
				url: `${API.si.tools}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Tools' }],
		}),
	}),
})

export const {
	useGetToolsMenuQuery,
	useToggleFavoriteMutation,
	useCreateToolsMenuMutation,
	useUpdateToolsMenuMutation,
	useDeleteToolsMenuMutation,
} = toolsMenuApiSlice
