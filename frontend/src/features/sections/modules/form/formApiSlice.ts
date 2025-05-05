import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { ICreateFormFieldDTO, ICreateFormStep } from './types/create'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const formApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getCreateFormSteps: builder.query<{ data: ICreateFormStep[] }, string>({
			query: section => ({
				url: API.createForm,
				params: new URLSearchParams({ section }),
			}),
			providesTags: [{ type: 'CreateForm', id: 'List' }],
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
		createFieldToCreateForm: builder.mutation<null, ICreateFormFieldDTO>({
			query: body => ({
				url: API.createForm,
				method: 'POST',
				body,
			}),
			invalidatesTags: [{ type: 'CreateForm', id: 'List' }],
		}),
		updateFieldToCreateForm: builder.mutation<null, ICreateFormFieldDTO>({
			query: body => ({
				url: `${API.createForm}/${body.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [{ type: 'CreateForm', id: 'List' }],
		}),
		updateSeveralFieldsToCreateForm: builder.mutation<null, ICreateFormFieldDTO[]>({
			query: body => ({
				url: `${API.createForm}/several`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [{ type: 'CreateForm', id: 'List' }],
		}),
		deleteFieldToCreateForm: builder.mutation<null, string>({
			query: id => ({
				url: `${API.createForm}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [{ type: 'CreateForm', id: 'List' }],
		}),
	}),
})

export const {
	useGetCreateFormStepsQuery,
	useCreateFieldToCreateFormMutation,
	useUpdateFieldToCreateFormMutation,
	useUpdateSeveralFieldsToCreateFormMutation,
	useDeleteFieldToCreateFormMutation,
} = formApiSlice
