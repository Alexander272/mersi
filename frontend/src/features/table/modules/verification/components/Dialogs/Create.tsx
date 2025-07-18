import { FC } from 'react'

import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { changeDialogIsOpen, getDialogState } from '@/features/dialog/dialogSlice'
import { Dialog } from '@/features/dialog/components/Dialog'
import { Create } from '../Forms/Create'

type Context = string | string[]

type Props = {
	title: string
}

export const CreateVerificationDialog: FC<Props> = ({ title }) => {
	const modal = useAppSelector(getDialogState('NewVerification'))
	const dispatch = useAppDispatch()

	const closeHandler = () => {
		dispatch(changeDialogIsOpen({ variant: 'NewVerification', isOpen: false }))
	}

	const context = modal?.context as Context
	return (
		<Dialog
			title={title}
			body={<Create ids={typeof context == 'string' ? [context] : context} />}
			open={modal?.isOpen || false}
			onClose={closeHandler}
			maxWidth='md'
			fullWidth
		/>
	)
}
