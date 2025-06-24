import { getPermissions } from '@/features/user/userSlice'
import { useAppSelector } from '@/hooks/redux'

export const useCheckPermission = (rule: string) => {
	const permissions = useAppSelector(getPermissions)
	if (!permissions.length) return false

	for (let i = 0; i < permissions.length; i++) {
		if (permissions[i] === rule) return true
	}
	return false
}
