import type { IFilter } from '@/features/table/types/table'

export interface IUser {
	id: string
	name: string
	role: string
	permissions: string[]
	token: string
	filters: IFilter[]
}

export interface IUserData {
	id: string
	ssoId: string
	username: string
	firstName: string
	lastName: string
	email: string
}
