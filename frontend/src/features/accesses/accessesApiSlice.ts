import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IAccesses, IAccessesDTO } from './types/accesses'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const accessesApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getAccesses: builder.query<{ data: IAccesses[] }, string>({
			query: realm => ({
				url: API.accesses,
				method: 'GET',
				params: { realm },
			}),
			providesTags: [{ type: 'Accesses', id: 'ALL' }],
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

		createAccess: builder.mutation<null, IAccessesDTO>({
			query: data => ({
				url: API.accesses,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'Accesses', id: 'ALL' },
				{ type: 'Users', id: arg.realmId },
			],
		}),
		updateAccess: builder.mutation<null, IAccessesDTO>({
			query: data => ({
				url: `${API.accesses}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'Accesses', id: 'ALL' }],
		}),
		deleteAccess: builder.mutation<null, string>({
			query: id => ({
				url: `${API.accesses}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'Accesses', id: 'ALL' },
				{ type: 'Users', id: arg },
			],
		}),
	}),
})

export const { useGetAccessesQuery, useCreateAccessMutation, useUpdateAccessMutation, useDeleteAccessMutation } =
	accessesApiSlice
