import { useCallback, useState } from 'react'

import type { ISI } from '../types/si'

export type Coordinates = {
	mouseX: number
	mouseY: number
}

// export type Status = 'reserve' | 'used' | 'moved'

export const useContextMenu = () => {
	const [coordinates, setCoordinates] = useState<Coordinates>()
	const [item, setItem] = useState<ISI>()
	// const [status, setStatus] = useState<Status>('reserve')
	const [isSelected, setIsSelected] = useState(false)

	const positionHandler = useCallback((coordinates?: Coordinates, item?: ISI, isSelected: boolean = false) => {
		setCoordinates(coordinates)
		setItem(item)
		setIsSelected(isSelected)
	}, [])

	return { coordinates, item, isSelected, positionHandler }
}
