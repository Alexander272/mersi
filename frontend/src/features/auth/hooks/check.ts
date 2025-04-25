import { getMenu } from '@/features/user/userSlice'
import { useAppSelector } from '@/hooks/redux'

export const useCheckPermission = (rule: string) => {
	const menu = useAppSelector(getMenu)
	if (!menu.length) return false

	for (let i = 0; i < menu.length; i++) {
		if (menu[i] === rule) return true
	}
	return false
}
