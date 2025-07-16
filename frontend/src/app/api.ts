export const API = {
	auth: {
		signIn: `auth/sign-in` as const,
		refresh: `auth/refresh` as const,
		signOut: `auth/sign-out` as const,
	},
	si: {
		base: 'si' as const,
		position: 'si/position' as const,
		moved: 'si/moved' as const,
		save: 'si/save' as const,
		instruments: {
			base: 'si/instruments' as const,
			unique: 'si/instruments/unique' as const,
		},
		documents: {
			base: 'si/documents' as const,
			list: 'si/documents/list' as const,
		},

		verification: {
			base: 'si/verifications' as const,
			all: 'si/verifications/all' as const,
			last: 'si/verifications/last' as const,
			fields: 'si/verifications/fields' as const,
		},
		location: {
			base: 'si/locations' as const,
			all: 'si/locations/all' as const,
		},
		export: 'files' as const,
		schedule: 'files/schedule' as const,
		context: 'context-menu' as const,
		tools: 'tools-menu' as const,
		repair: 'repair' as const,
		preservation: 'preservation' as const,
		transferToSave: 'transfer-to-save' as const,
		transferToDep: 'transfer-to-department' as const,
	},
	departments: '/departments' as const,
	employees: '/employees' as const,
	responsible: '/responsible' as const,
	channels: '/channels' as const,
	filters: '/filters' as const,
	users: {
		base: '/users' as const,
		sync: '/users/sync' as const,
		access: '/users/access' as const,
		realm: '/users/realm' as const,
	},
	roles: '/roles' as const,
	realms: {
		base: '/realms' as const,
		user: '/realms/user' as const,
		choose: '/realms/choose' as const,
	},
	accesses: '/accesses' as const,
	sections: {
		base: '/sections' as const,
		grouped: '/sections/grouped' as const,
	},
	columns: '/columns' as const,
	createForm: '/forms/create',
}
