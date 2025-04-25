import { PayloadAction, createSlice } from '@reduxjs/toolkit'

import type { RootState } from '@/app/store'
import type { IUser } from './types/user'

interface IUserState {
	// user?: IUser

	id: string | null
	name?: string
	role: string | null
	menu: string[]
	token: string | null
}

const initialState: IUserState = {
	id: null,
	role: null,
	token: null,
	menu: [],
}

const userSlice = createSlice({
	name: 'user',
	initialState,
	reducers: {
		setUser: (state, action: PayloadAction<IUser>) => {
			// state.user = action.payload

			state.id = action.payload.id
			state.name = action.payload.name
			state.role = action.payload.role
			state.menu = action.payload.menu
			state.token = action.payload.token
		},

		setRole: (state, action: PayloadAction<string>) => {
			state.role = action.payload
		},
		setMenu: (state, action: PayloadAction<string[]>) => {
			state.menu = action.payload
		},

		resetUser: () => initialState,
	},
})

export const getToken = (state: RootState) => state.user.token
export const getMenu = (state: RootState) => state.user.menu

export const userPath = userSlice.name
export const userReducer = userSlice.reducer

export const { setUser, setRole, setMenu, resetUser } = userSlice.actions
