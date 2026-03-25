import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { IGetSiDTO } from '../../types/si'
import type { IPeriod } from './types/period'
import { API } from '@/app/api'
import { HttpCodes } from '@/constants/httpCodes'
import { buildSiUrlParams } from '../../utils/buildUrlParams'
import { saveAs } from '@/features/files/utils/saveAs'
import { apiSlice } from '@/app/apiSlice'

const exportApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		export: builder.query<null, IGetSiDTO>({
			queryFn: async (params, _api, _, baseQuery) => {
				const filename = `Список инструментов от ${dayjs().format('DD-MM-YYYY')}.xlsx`
				const result = await baseQuery({
					url: API.si.export,
					params: buildSiUrlParams(params),
					cache: 'no-cache',
					responseHandler: response => (response.status === HttpCodes.OK ? response.blob() : response.json()),
				})

				if (result.error) {
					console.log(result.error)
					const fetchError = result.error as IFetchError
					toast.error(fetchError.data.message, { autoClose: false })
				}

				if (result.data instanceof Blob) saveAs(result.data, filename)
				return { data: null }
			},
		}),

		makeScheduler: builder.query<null, IPeriod>({
			queryFn: async (params, _api, _, baseQuery) => {
				const filename = `График поверки от ${dayjs().format('DD-MM-YYYY')}.xlsx`
				const result = await baseQuery({
					url: API.si.schedule,
					params: new URLSearchParams({
						'period[gte]': params.gte.toString(),
						'period[lte]': params.lte.toString(),
						section: params.section,
					}),
					cache: 'no-cache',
					responseHandler: response => (response.status === HttpCodes.OK ? response.blob() : response.json()),
				})

				if (result.error) {
					console.log(result.error)
					const fetchError = result.error as IFetchError
					toast.error(fetchError.data.message, { autoClose: false })
				}

				if (result.data instanceof Blob) saveAs(result.data, filename)
				return { data: null }
			},
		}),

		makeAccountingLog: builder.query<null, IPeriod>({
			queryFn: async (params, _api, _, baseQuery) => {
				const filename = `Журнал учета средств измерения от ${dayjs().format('DD-MM-YYYY')}.xlsx`
				const result = await baseQuery({
					url: API.si.log,
					params: new URLSearchParams({
						section: params.section,
					}),
					cache: 'no-cache',
					responseHandler: response => (response.status === HttpCodes.OK ? response.blob() : response.json()),
				})

				if (result.error) {
					console.log(result.error)
					const fetchError = result.error as IFetchError
					toast.error(fetchError.data.message, { autoClose: false })
				}

				if (result.data instanceof Blob) saveAs(result.data, filename)
				return { data: null }
			},
		}),
	}),
})

export const { useMakeSchedulerQuery, useLazyMakeSchedulerQuery, useLazyMakeAccountingLogQuery, useLazyExportQuery } =
	exportApiSlice
