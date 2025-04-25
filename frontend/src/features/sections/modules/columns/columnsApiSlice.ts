import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IColumn, IColumnDTO } from './types/columns'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const columnsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getColumns: builder.query<{ data: IColumn[] }, string>({
			query: section => ({
				url: API.columns,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Columns', id: 'List' }],
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
		createColumn: builder.mutation<null, IColumnDTO>({
			query: body => ({
				url: API.columns,
				method: 'POST',
				body,
			}),
			invalidatesTags: [{ type: 'Columns', id: 'List' }],
		}),
		updateColumn: builder.mutation<null, IColumnDTO>({
			query: body => ({
				url: `${API.columns}/${body.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [{ type: 'Columns', id: 'List' }],
		}),
		deleteColumn: builder.mutation<null, string>({
			query: id => ({
				url: `${API.columns}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [{ type: 'Columns', id: 'List' }],
		}),
	}),
})

export const { useGetColumnsQuery, useCreateColumnMutation, useUpdateColumnMutation, useDeleteColumnMutation } =
	columnsApiSlice
