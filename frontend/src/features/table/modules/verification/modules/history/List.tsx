import { FC } from 'react'

import { VerificationListBase } from '../../VerificationListBase'
import { VerificationHistoryDialog } from './Dialog'
import { ListItem } from './Item'

type Props = {
	section: string
}

export const VerificationHistoryList: FC<Props> = ({ section }) => {
	return (
		<VerificationListBase
			section={section}
			group='history'
			dialogVariant='EditVerificationHistory'
			toastMessage='Колонки сохранены'
			DialogComponent={VerificationHistoryDialog}
			ItemComponent={ListItem}
		/>
	)
}
