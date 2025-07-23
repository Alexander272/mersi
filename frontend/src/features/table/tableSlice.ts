import { type PayloadAction, createSlice } from '@reduxjs/toolkit'

import type { IContextMenu, ISelect } from './types/table'
import type { IFilter, ISearch, ISort } from './types/params'
import type { Status } from './types/si'
import { Size } from '@/constants/defaultValues'
import { RootState } from '@/app/store'
import { localKeys } from './constants/storage'

interface ITableSlice {
	page: number
	size: number
	status: Status
	sort: ISort
	filters: IFilter[]
	search: ISearch
	selected: ISelect
	contextMenu?: IContextMenu
	hidden: ISelect
}

const initialState: ITableSlice = {
	page: +(localStorage.getItem(localKeys.page) || 1),
	size: +(localStorage.getItem(localKeys.size) || Size),
	status: 'work',
	sort: {},
	filters: [],
	search: {
		value: '',
		fields: ['name', 'factoryNumber'],
	},
	selected: {},
	hidden: JSON.parse(localStorage.getItem(localKeys.hidden) || '{}'),
}

const tableSlice = createSlice({
	name: 'table',
	initialState,
	reducers: {
		setPage: (state, action: PayloadAction<number>) => {
			state.page = action.payload
			localStorage.setItem(localKeys.page, action.payload.toString())
		},
		setSize: (state, action: PayloadAction<number>) => {
			state.size = action.payload
			localStorage.setItem(localKeys.size, action.payload.toString())
		},
		setStatus: (state, action: PayloadAction<Status>) => {
			state.status = action.payload
		},

		setSort: (state, action: PayloadAction<string>) => {
			if (!state.sort[action.payload]) {
				state.sort = { ...(state.sort || {}), [action.payload]: 'ASC' }
				return
			}

			if (state.sort[action.payload] == 'ASC') state.sort[action.payload] = 'DESC'
			else {
				delete state.sort[action.payload]
			}
		},

		setFilters: (state, action: PayloadAction<IFilter[]>) => {
			state.filters = action.payload
		},

		setSearch: (state, action: PayloadAction<string>) => {
			state.search.value = action.payload
		},
		setSearchFields: (state, action: PayloadAction<string[]>) => {
			state.search.fields = action.payload
		},

		setSelected: (state, action: PayloadAction<string | string[] | undefined>) => {
			if (action.payload) {
				if (typeof action.payload == 'string') {
					if (state.selected[action.payload]) delete state.selected[action.payload]
					else state.selected[action.payload] = true
				} else {
					state.selected = action.payload.reduce((a, v) => ({ ...a, [v]: true }), {})
				}
			} else state.selected = {}
		},

		setContextMenu: (state, action: PayloadAction<IContextMenu | undefined>) => {
			state.contextMenu = action.payload
		},

		setHidden: (state, action: PayloadAction<ISelect | undefined>) => {
			if (action.payload) state.hidden = action.payload
			else state.hidden = {}
			localStorage.setItem(localKeys.hidden, JSON.stringify(state.hidden))
		},

		resetTable: () => {
			localStorage.removeItem(localKeys.page)
			return initialState
		},
	},
})

export const getTablePage = (state: RootState) => state.table.page
export const getTableSize = (state: RootState) => state.table.size
export const getStatus = (state: RootState) => state.table.status
export const getSort = (state: RootState) => state.table.sort
export const getFilters = (state: RootState) => state.table.filters
export const getSearch = (state: RootState) => state.table.search
export const getSelected = (state: RootState) => state.table.selected
export const getContextMenu = (state: RootState) => state.table.contextMenu
export const getHidden = (state: RootState) => state.table.hidden

export const tablePath = tableSlice.name
export const tableReducer = tableSlice.reducer

export const {
	setPage,
	setSize,
	setStatus,
	setSort,
	setFilters,
	setSearch,
	setSearchFields,
	setSelected,
	setContextMenu,
	setHidden,
	resetTable,
} = tableSlice.actions
