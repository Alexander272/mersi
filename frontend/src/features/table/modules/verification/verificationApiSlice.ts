import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IVerificationField } from './types/verificationFields'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { IVerificationDTO } from './types/verification'

const verificationApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getVerificationFields: builder.query<{ data: IVerificationField[] }, string>({
			query: section => ({
				url: `${API.si.verification.base}/fields`,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'Verification', id: 'Fields' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createVerification: builder.mutation<null, IVerificationDTO>({
			query: body => ({
				url: API.si.verification.base,
				method: 'POST',
				body,
			}),
			invalidatesTags: [{ type: 'SI', id: 'ALL' }],
		}),
	}),
})

export const { useGetVerificationFieldsQuery, useCreateVerificationMutation } = verificationApiSlice
