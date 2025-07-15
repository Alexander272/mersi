import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { ITransferToSave, ITransferToSaveDTO } from './types/save'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const transferApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getTransferToSave: builder.query<{ data: ITransferToSave[] }, string>({
			query: instrument => ({
				url: API.si.transferToSave,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'TransferToSave' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getLastTransferToSave: builder.query<{ data: ITransferToSave }, string>({
			query: instrument => ({
				url: `${API.si.transferToSave}/last`,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'LastTransferToSave' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					if (fetchError.status == 404) return
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createTransferToSave: builder.mutation<null, ITransferToSaveDTO>({
			query: body => ({
				url: API.si.transferToSave,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'TransferToSave' },
				{ type: 'SI', id: 'LastTransferToSave' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
		updateTransferToSave: builder.mutation<null, ITransferToSaveDTO>({
			query: body => ({
				url: `${API.si.transferToSave}/${body.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'TransferToSave' },
				{ type: 'SI', id: 'LastTransferToSave' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const {
	useGetTransferToSaveQuery,
	useGetLastTransferToSaveQuery,
	useCreateTransferToSaveMutation,
	useUpdateTransferToSaveMutation,
} = transferApiSlice
