import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IFilter, ISortDTO } from './types/params'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const filterApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getFilters: builder.query<{ data: IFilter[] }, string>({
			query: section => ({
				url: API.filters,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Sections', id: 'Filters' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		saveFilters: builder.mutation<void, { filters: IFilter[]; section: string }>({
			query: data => ({
				url: `${API.filters}/change`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Filters' }],
		}),

		getSort: builder.query<{ data: ISortDTO[] }, string>({
			query: section => ({
				url: API.sorting,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Sections', id: 'Sort' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		saveSort: builder.mutation<void, { sort: ISortDTO[]; section: string }>({
			query: data => ({
				url: `${API.sorting}/change`,
				method: 'POST',
				params: new URLSearchParams({ section: data.section }),
				body: data.sort,
			}),
			invalidatesTags: [{ type: 'Sections', id: 'Sort' }],
		}),
	}),
})

export const { useGetFiltersQuery, useSaveFiltersMutation, useGetSortQuery, useSaveSortMutation } = filterApiSlice
