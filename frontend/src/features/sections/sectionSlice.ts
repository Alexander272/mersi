import { PayloadAction, createSlice } from '@reduxjs/toolkit'

import type { RootState } from '@/app/store'
import type { ISection } from './types/sections'
import { SectionKey } from './constants/storage'

interface ISectionState {
	section: ISection | null
}

const initialState: ISectionState = {
	section: JSON.parse(localStorage.getItem(SectionKey) || 'null'),
}

const sectionSlice = createSlice({
	name: 'section',
	initialState,
	reducers: {
		setSection: (state, action: PayloadAction<ISection>) => {
			state.section = action.payload
			localStorage.setItem(SectionKey, JSON.stringify(state.section))
		},

		resetSection: () => initialState,
	},
})

export const getSection = (state: RootState) => state.section.section

export const sectionPath = sectionSlice.name
export const sectionReducer = sectionSlice.reducer

export const { setSection, resetSection } = sectionSlice.actions
