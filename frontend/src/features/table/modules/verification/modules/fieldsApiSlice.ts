import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IVerificationField, IVerificationFieldDTO } from './types/verificationFields'
import { API } from '@/app/api'
import { HttpCodes } from '@/constants/httpCodes'
import { apiSlice } from '@/app/apiSlice'

const verificationFieldsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getVerificationFields: builder.query<{ data: IVerificationField[] }, { section: string; group: string }>({
			query: req => ({
				url: API.si.verification.fields,
				params: new URLSearchParams({ section: req.section, group: req.group }),
			}),
			providesTags: [{ type: 'Verification', id: 'Fields' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					if (fetchError.status == HttpCodes.NOT_FOUND) return
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createVerificationFields: builder.mutation<null, IVerificationFieldDTO[]>({
			query: data => ({
				url: `${API.si.verification.fields}/several`,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [{ type: 'Verification', id: 'Fields' }],
		}),
		updateVerificationFields: builder.mutation<null, IVerificationFieldDTO[]>({
			query: data => ({
				url: `${API.si.verification.fields}/several`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [{ type: 'Verification', id: 'Fields' }],
		}),
		deleteVerificationFields: builder.mutation<null, string[]>({
			query: data => ({
				url: `${API.si.verification.fields}/several`,
				method: 'DELETE',
				body: data,
			}),
			invalidatesTags: [{ type: 'Verification', id: 'Fields' }],
		}),
	}),
})

export const {
	useGetVerificationFieldsQuery,
	useCreateVerificationFieldsMutation,
	useUpdateVerificationFieldsMutation,
	useDeleteVerificationFieldsMutation,
} = verificationFieldsApiSlice
