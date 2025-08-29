import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IGroupedSections, ISection, ISectionDTO } from './types/sections'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const sectionsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getSections: builder.query<{ data: ISection[] }, string>({
			query: realm => ({
				url: API.sections.base,
				params: new URLSearchParams({ realm }),
			}),
			providesTags: [{ type: 'Sections', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					console.log(error)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getGroupedSections: builder.query<{ data: IGroupedSections[] }, null>({
			query: () => API.sections.grouped,
			providesTags: [{ type: 'Sections', id: 'Grouped' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					console.log(error)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createSection: builder.mutation<{ id: string }, ISectionDTO>({
			query: section => ({
				url: API.sections.base,
				method: 'POST',
				body: section,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'ALL' }],
		}),

		updateSection: builder.mutation<null, ISectionDTO>({
			query: section => ({
				url: `${API.sections.base}/${section.id}`,
				method: 'PUT',
				body: section,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'ALL' }],
		}),

		deleteSection: builder.mutation<null, string>({
			query: id => ({
				url: `${API.sections.base}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [{ type: 'Sections', id: 'ALL' }],
		}),
	}),
})

export const {
	useGetSectionsQuery,
	useGetGroupedSectionsQuery,
	useCreateSectionMutation,
	useUpdateSectionMutation,
	useDeleteSectionMutation,
} = sectionsApiSlice
