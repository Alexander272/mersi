export const PermRules = Object.freeze({
	SI: {
		Read: 'si:read' as const,
		Write: 'si:write' as const,
	},
	Location: {
		Read: 'location:read' as const,
		Write: 'location:write' as const,
	},
	Verification: {
		Read: 'verification:read' as const,
		Write: 'verification:write' as const,
	},
	Documents: {
		Write: 'documents:write' as const,
	},
	Employee: {
		Read: 'employee:read' as const,
		Write: 'employee:write' as const,
	},
	Department: {
		Read: 'department:read' as const,
		Write: 'department:write' as const,
	},
	Reserve: {
		// Read: 'reserve:read' as const,
		Write: 'reserve:write' as const,
	},
	Realms: {
		Read: 'realms:read' as const,
		Write: 'realms:write' as const,
	},
})
