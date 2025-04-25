import { JSX } from 'react'

export const formatNumber = (num: number | undefined): string | JSX.Element => {
	if (!num) return '0'

	if (Math.abs(num) < 0.0001) {
		let tmp = num.toLocaleString('ru-RU', { notation: 'scientific' })
		tmp = tmp.replace(/E|e/g, '*10^')
		const parts = tmp.split('^')

		return (
			<>
				{parts[0]}
				<sup>{parts[1]}</sup>
			</>
		)
	} else {
		return num.toLocaleString('ru-Ru')
	}
}

export const numberFormat = new Intl.NumberFormat('ru-Ru').format

export const removeSpace = (values: string[]) => {
	return values.map(v => v.replace(/\s+/g, ''))
}
