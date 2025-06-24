import { PermRules } from '@/constants/permissions'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { getSection } from '@/features/sections/sectionSlice'
import { useGetSIQuery } from '../siApiSlice'
import { getFilters, getSort, getTablePage, getTableSize } from '../tableSlice'

export const useGetSI = () => {
	// const status = useAppSelector(getStatus)
	const section = useAppSelector(getSection)
	const page = useAppSelector(getTablePage)
	const size = useAppSelector(getTableSize)

	const sort = useAppSelector(getSort)
	const filters = useAppSelector(getFilters)

	const all = useCheckPermission(PermRules.Location.Write)

	const query = useGetSIQuery(
		// { section: section?.id || '', status, page, size, all, sort, filter },
		{ section: section?.id || '', page, size, all, sort, filters },
		{ skip: !section?.id, pollingInterval: 5 * 60000, skipPollingIfUnfocused: true /*refetchOnFocus: true*/ }
	)

	return query
}
