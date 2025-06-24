import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IGroupedSections, ISection } from './types/sections'
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
	}),
})

export const { useGetSectionsQuery, useGetGroupedSectionsQuery } = sectionsApiSlice
