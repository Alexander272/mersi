import { FC } from 'react'

import { DocIcon } from '@/components/Icons/DocIcon'
import { ImageIcon } from '@/components/Icons/ImageIcon'
import { PdfIcon } from '@/components/Icons/PdfIcon'
import { SheetIcon } from '@/components/Icons/SheetIcon'

type Props = {
	type: string
}

export const FileIcon: FC<Props> = ({ type }) => {
	switch (type) {
		case 'pdf':
			return <PdfIcon />
		case 'image':
			return <ImageIcon />
		case 'sheet':
			return <SheetIcon />
		default:
			return <DocIcon />
	}
}
