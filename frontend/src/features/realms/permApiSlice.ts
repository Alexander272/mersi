import { apiSlice } from '@/app/apiSlice'

export const permApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		reloadPerms: builder.mutation<void, null>({
			query: () => ({
				url: '/permissions/reload',
				method: 'POST',
			}),
		}),
	}),
})

export const { useReloadPermsMutation } = permApiSlice
