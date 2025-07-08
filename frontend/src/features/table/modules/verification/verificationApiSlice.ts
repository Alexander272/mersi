import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IVerificationField } from './types/verificationFields'
import type { IVerification, IVerificationDTO } from './types/verification'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const verificationApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getVerificationFields: builder.query<{ data: IVerificationField[] }, string>({
			query: section => ({
				url: API.si.verification.fields,
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

		getLastVerification: builder.query<{ data: IVerification }, string>({
			query: instrument => ({
				url: API.si.verification.last,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'Verification', id: 'Last' }],
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
			invalidatesTags: [
				{ type: 'SI', id: 'ALL' },
				{ type: 'Verification', id: 'Last' },
			],
		}),
	}),
})

export const { useGetVerificationFieldsQuery, useGetLastVerificationQuery, useCreateVerificationMutation } =
	verificationApiSlice
