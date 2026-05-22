import { FC } from 'react'

import { VerificationListBase } from '../../VerificationListBase'
import { VerificationFieldsDialog } from './Dialog'
import { FieldItem } from './Item'

type Props = {
	section: string
}

export const VerificationFieldsList: FC<Props> = ({ section }) => {
	return (
		<VerificationListBase
			section={section}
			group='form'
			dialogVariant='EditVerificationFields'
			toastMessage='Поля формы сохранены'
			DialogComponent={VerificationFieldsDialog}
			ItemComponent={FieldItem}
		/>
	)
}
