import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IForcedReceipt, ILocation, ILocationDTO, IReceiving } from './types/location'
import { API } from '@/app/api'
import { HttpCodes } from '@/constants/httpCodes'
import { apiSlice } from '@/app/apiSlice'

const locationApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getLocations: builder.query<{ data: ILocation[] }, string>({
			query: instrument => ({
				url: API.si.locations.base,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'Location', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getLastLocation: builder.query<{ data: ILocation }, string>({
			query: instrument => ({
				url: API.si.locations.last,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'Location', id: 'LAST' }],
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

		createLocation: builder.mutation<null, ILocationDTO>({
			query: data => ({
				url: API.si.locations.base,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
			],
		}),
		createSeveralLocations: builder.mutation<{ message: string }, ILocationDTO[]>({
			query: data => ({
				url: API.si.locations.several,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
			],
		}),

		updateLocation: builder.mutation<null, ILocationDTO>({
			query: data => ({
				url: `${API.si.locations.base}/${data.id}`,
				method: 'PUT',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
			],
		}),

		receiving: builder.mutation<null, IReceiving>({
			query: data => ({
				url: API.si.locations.receiving,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
			],
		}),
		forcedReceiving: builder.mutation<null, IForcedReceipt>({
			query: data => ({
				url: API.si.locations.forced,
				method: 'POST',
				body: data,
			}),
			invalidatesTags: [
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
			],
		}),

		deleteLocation: builder.mutation<null, string>({
			query: id => ({
				url: `${API.si.locations.base}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: [
				{ type: 'Location', id: 'ALL' },
				{ type: 'Location', id: 'LAST' },
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'MOVED' },
			],
		}),
	}),
})

export const {
	useGetLocationsQuery,
	useGetLastLocationQuery,
	useCreateLocationMutation,
	useCreateSeveralLocationsMutation,
	useUpdateLocationMutation,
	useReceivingMutation,
	useForcedReceivingMutation,
	useDeleteLocationMutation,
} = locationApiSlice
