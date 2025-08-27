import { toast } from 'react-toastify'

import type { IBaseFetchError, IFetchError } from '@/app/types/error'
import type { IDocument } from './types/document'
import type { DeleteDocuments, IGetDocumentsDTO, IUploadFiles } from './types/request'
import { API } from '@/app/api'
import { HttpCodes } from '@/constants/httpCodes'
import { apiSlice } from '@/app/apiSlice'
import { saveAs } from './utils/saveAs'

const filesApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getTempFiles: builder.query<{ data: IDocument[] }, IGetDocumentsDTO>({
			query: req => ({
				url: `${API.si.documents.temp}/${req.group}`,
				method: 'GET',
				params: new URLSearchParams({ instrument: req.instrument }),
			}),
			providesTags: [{ type: 'Documents', id: 'Temp' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getFileList: builder.query<{ data: IDocument[] }, IGetDocumentsDTO>({
			query: req => ({
				url: `${API.si.documents.list}/${req.group}`,
				method: 'GET',
				params: new URLSearchParams({ instrument: req.instrument }),
			}),
			providesTags: [{ type: 'Documents', id: 'List' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		downloadFile: builder.query<null, { path: string; label: string }>({
			queryFn: async (doc, _api, _, baseQuery) => {
				const result = await baseQuery({
					url: API.si.documents.base,
					params: new URLSearchParams({ path: doc.path }),
					cache: 'no-cache',
					responseHandler: response => (response.status === HttpCodes.OK ? response.blob() : response.json()),
				})

				if (result.error) {
					console.log(result.error)
					const fetchError = result.error as IFetchError
					toast.error(fetchError.data.message, { autoClose: false })
				}

				if (result.data instanceof Blob) saveAs(result.data, doc.label)
				return { data: null }
			},
		}),
		uploadFiles: builder.mutation<{ data: IDocument[] }, IUploadFiles>({
			query: data => ({
				url: `${API.si.documents.base}`,
				method: 'POST',
				body: data.data,
				validateStatus: response => response.status === HttpCodes.CREATED,
			}),
			invalidatesTags: [
				{ type: 'Documents', id: 'List' },
				{ type: 'Documents', id: 'Temp' },
			],
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
		deleteFile: builder.mutation<null, DeleteDocuments>({
			query: data => ({
				url: `${API.si.documents.base}/${data.id}`,
				method: 'DELETE',
				params: new URLSearchParams({
					instrumentId: data.instrumentId,
					group: data.group,
					filename: data.filename,
					isTemp: data.isTemp.toString(),
				}),
			}),
			invalidatesTags: [
				{ type: 'Documents', id: 'List' },
				{ type: 'Documents', id: 'Temp' },
			],
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
	}),
})

export const {
	useGetTempFilesQuery,
	useGetFileListQuery,
	useUploadFilesMutation,
	useLazyDownloadFileQuery,
	useDeleteFileMutation,
} = filesApiSlice
