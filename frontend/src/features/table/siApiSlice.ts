import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IChangePositionDTO, IGetSiDTO, ISI, ISiDTO } from './types/si'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { buildSiUrlParams } from './utils/buildUrlParams'

const SIApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getSI: builder.query<{ data: ISI[]; total: number }, IGetSiDTO>({
			query: params => ({
				url: API.si.base,
				method: 'GET',
				params: buildSiUrlParams(params),
			}),
			providesTags: [{ type: 'SI', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					console.error(fetchError)
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getSIById: builder.query<{ data: ISiDTO }, string>({
			query: id => `${API.si.base}/${id}`,
			providesTags: (_res, _err, arg) => [
				{ type: 'SI', id: arg },
				{ type: 'SI', id: 'ID' },
			],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					console.error(fetchError)
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createSi: builder.mutation<null, ISiDTO>({
			query: body => ({
				url: API.si.base,
				method: 'POST',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'ALL' },
				{ type: 'Instrument', id: 'Unique' },
			],
		}),

		updateSI: builder.mutation<null, ISiDTO>({
			query: body => ({
				url: `${API.si.base}/${body.instrument.id}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: arg.instrument.id },
				{ type: 'Instrument', id: 'Unique' },
			],
		}),

		changePosition: builder.mutation<null, IChangePositionDTO>({
			query: body => ({
				url: `${API.si.position}`,
				method: 'PUT',
				body,
			}),
			invalidatesTags: [
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: 'ID' },
			],
		}),

		deleteSI: builder.mutation<null, string>({
			query: id => ({
				url: `${API.si.base}/${id}`,
				method: 'DELETE',
			}),
			invalidatesTags: (_res, _err, arg) => [
				{ type: 'SI', id: 'ALL' },
				{ type: 'SI', id: arg },
				{ type: 'Instrument', id: 'Unique' },
			],
		}),
	}),
})

export const {
	useGetSIQuery,
	useGetSIByIdQuery,
	useCreateSiMutation,
	useUpdateSIMutation,
	useChangePositionMutation,
	useDeleteSIMutation,
} = SIApiSlice
