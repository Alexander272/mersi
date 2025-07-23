import { PermRules } from '@/constants/permissions'
import { useAppSelector } from '@/hooks/redux'
import { useCheckPermission } from '@/features/user/hooks/check'
import { useGetSIQuery } from '../siApiSlice'
import { getSection } from '@/features/sections/sectionSlice'
import { getFilters, getSearch, getSort, getStatus, getTablePage, getTableSize } from '../tableSlice'

export const useGetSI = () => {
	const status = useAppSelector(getStatus)
	const section = useAppSelector(getSection)
	const page = useAppSelector(getTablePage)
	const size = useAppSelector(getTableSize)

	const search = useAppSelector(getSearch)
	const sort = useAppSelector(getSort)
	const filters = useAppSelector(getFilters)

	const all = useCheckPermission(PermRules.Location.Write)

	const query = useGetSIQuery(
		{ section: section?.id || '', page, size, status, all, sort, filters, search },
		{ skip: !section?.id, pollingInterval: 5 * 60000, skipPollingIfUnfocused: true /*refetchOnFocus: true*/ }
	)

	return query
}
