import type { ISiDTO } from './types/si'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

const SIApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
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
	}),
})

export const { useCreateSiMutation } = SIApiSlice
