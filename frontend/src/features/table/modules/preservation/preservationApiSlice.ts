import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IPreservation, IPreservationDTO } from './types/preservation'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const preservationApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getPreservations: builder.query<{ data: IPreservation[] }, string>({
			query: instrument => ({
				url: API.si.preservation,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'Preservation' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getLastPreservation: builder.query<{ data: IPreservation }, string>({
			query: instrument => ({
				url: `${API.si.preservation}/last`,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'LastPreservation' }],
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

		createPreservation: builder.mutation<null, IPreservationDTO>({
			query: body => ({
				url: API.si.preservation,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'Preservation' },
				{ type: 'SI', id: 'LastPreservation' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
		updatePreservation: builder.mutation<null, IPreservationDTO>({
			query: body => ({
				url: `${API.si.preservation}/${body.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'Preservation' },
				{ type: 'SI', id: 'LastPreservation' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const {
	useGetPreservationsQuery,
	useGetLastPreservationQuery,
	useCreatePreservationMutation,
	useUpdatePreservationMutation,
} = preservationApiSlice
