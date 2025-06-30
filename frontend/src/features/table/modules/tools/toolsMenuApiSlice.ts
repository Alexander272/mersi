import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IToggleToolsMenu, IToolsMenu } from './types/toolsMenu'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const toolsMenuApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getToolsMenu: builder.query<{ data: IToolsMenu[] }, string>({
			query: section => ({
				url: API.si.tools,
				params: new URLSearchParams({ section }),
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
	}),
})

export const { useGetToolsMenuQuery, useToggleFavoriteMutation } = toolsMenuApiSlice
