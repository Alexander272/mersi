import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IRepair, IRepairDTO } from './types/repair'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

const repairApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getRepair: builder.query<{ data: IRepair[] }, string>({
			query: instrument => ({
				url: API.si.repair,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'Repair' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getLastRepair: builder.query<{ data: IRepair }, string>({
			query: instrument => ({
				url: `${API.si.repair}/last`,
				params: new URLSearchParams({ instrument }),
			}),
			providesTags: [{ type: 'SI', id: 'Repair' }],
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

		createRepair: builder.mutation<null, IRepairDTO>({
			query: body => ({
				url: API.si.repair,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'Repair' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
		updateRepair: builder.mutation<null, IRepairDTO>({
			query: body => ({
				url: `${API.si.repair}/${body.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'Repair' },
				{ type: 'SI', id: 'ALL' },
			],
		}),
	}),
})

export const { useGetRepairQuery, useGetLastRepairQuery, useCreateRepairMutation, useUpdateRepairMutation } =
	repairApiSlice
