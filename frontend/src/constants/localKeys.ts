export const localKeys = Object.freeze({
	page: 'mersi/page' as const,
	size: 'mersi/size' as const,
	sort: 'mersi/sort' as const,
	filter: 'mersi/filter' as const,
	hidden: 'mersi/hidden' as const,
	columns: 'mersi/columns' as const,
	activeFilters: 'mersi/filters/active' as const,
	changedColumns: (key: string) => `mersi/changedColumns/${key}` as const,

	instrument: 'mersi/new/instrument' as const,
	verification: 'mersi/new/verification' as const,
	location: 'mersi/new/location' as const,
})

export const DraftKey = 'draft'
