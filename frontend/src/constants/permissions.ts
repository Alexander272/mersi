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
	Import: {
		Write: 'import:write' as const,
	},

	Repair: {
		Read: 'repair:read' as const,
		Write: 'repair:write' as const,
	},
	Preservation: {
		Read: 'preservation:read' as const,
		Write: 'preservation:write' as const,
	},
	TransferToDep: {
		Read: 'transfer-to-department:read' as const,
		Write: 'transfer-to-department:write' as const,
	},
	TransferToSave: {
		Read: 'transfer-to-save:read' as const,
		Write: 'transfer-to-save:write' as const,
	},
	WriteOff: {
		Read: 'write-off:read' as const,
		Write: 'write-off:write' as const,
	},
})
