import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IColumn, IColumnDTO, IColumnPositionDTO } from './types/columns'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { localKeys } from '@/constants/localKeys'

export const columnsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getColumns: builder.query<{ data: IColumn[] }, { section: string; original?: boolean }>({
			query: data => ({
				url: API.columns,
				params: new URLSearchParams({ section: data.section }),
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
			transformResponse(baseQueryReturnValue: { data: IColumn[] }, _meta, arg) {
				// console.log('baseQueryReturnValue', baseQueryReturnValue)
				// console.log('arg', arg)
				if (arg.original) return baseQueryReturnValue
				const changed = JSON.parse(localStorage.getItem(localKeys.changedColumns(arg.section)) || '{}')

				const newData = baseQueryReturnValue?.data.map(c => {
					if (!c.children?.length) return { ...c, ...changed?.[c.id] }
					const newChildren = c.children.map(c => ({ ...c, ...changed?.[c.id] }))
					return { ...c, ...changed?.[c.id], children: newChildren }
				})

				// console.log('newData', newData)

				return { data: newData }
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
		updateColumnPositions: builder.mutation<null, IColumnPositionDTO[]>({
			query: body => ({
				url: `${API.columns}/positions`,
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

export const {
	useGetColumnsQuery,
	useCreateColumnMutation,
	useUpdateColumnMutation,
	useUpdateColumnPositionsMutation,
	useDeleteColumnMutation,
} = columnsApiSlice
