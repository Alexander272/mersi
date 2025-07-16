import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { ITransferToDepartment, ITransferToDepartmentDTO } from './types/department'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const transferApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getTransferToDepartment: builder.query<{ data: ITransferToDepartment[] }, string>({
			query: instrument => ({
				url: API.si.transferToDep,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'TransferToDepartment' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createTransferToDepartment: builder.mutation<{ data: ITransferToDepartment }, ITransferToDepartmentDTO>({
			query: data => ({
				url: API.si.transferToDep,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'TransferToDepartment' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
		updateTransferToDepartment: builder.mutation<{ data: ITransferToDepartment }, ITransferToDepartmentDTO>({
			query: data => ({
				url: `${API.si.transferToDep}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'TransferToDepartment' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const {
	useGetTransferToDepartmentQuery,
	useCreateTransferToDepartmentMutation,
	useUpdateTransferToDepartmentMutation,
} = transferApiSlice
