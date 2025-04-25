export const API = {
	auth: {
		signIn: `auth/sign-in` as const,
		refresh: `auth/refresh` as const,
		signOut: `auth/sign-out` as const,
	},
	si: {
		instruments: {
			base: 'si/instruments' as const,
		},
		verification: {
			base: 'si/verifications' as const,
			all: 'si/verifications/all' as const,
		},
		location: {
			base: 'si/locations' as const,
			all: 'si/locations/all' as const,
		},
		documents: {
			base: 'si/verifications/documents' as const,
		},
		export: 'files' as const,
		schedule: 'files/schedule' as const,
		base: 'si' as const,
		moved: 'si/moved' as const,
		save: 'si/save' as const,
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
}
