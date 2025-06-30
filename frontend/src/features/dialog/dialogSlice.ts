import { RootState } from '@/app/store'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export type DialogVariants =
	| 'CreateTableItem'
	| 'EditTableItem'
	| 'ChangePosition'
	| 'NewVerification'
	| 'SeveralVerifications'
	| 'NewLocation'
	| 'SeveralLocations'
	| 'DeleteLocation'
	| 'SendToReserve'
	| 'Employee'
	| 'CreateDepartment'
	| 'EditDepartment'
	| 'Confirm'
	| 'ViewLocationHistory'
	| 'ViewVerificationHistory'
	| 'Period'
	| 'Documents'
	| 'Receive'
	| 'Access'
	| 'Columns'
	| 'CreateFormGroup'
	| 'CreateFormField'
	| 'AddRepair'
	| 'AddPreservation'
	| 'AddTransferToSave'
	| 'AddTransferToDep'
	| 'History'
	| 'WriteOff'

interface IDialogOptions {
	isOpen: boolean
	context?: unknown
}

type IDialogState = {
	[key in DialogVariants]?: IDialogOptions
}

interface IChangeDialogAction extends IDialogOptions {
	variant: DialogVariants
}

const initialState: IDialogState = {}

const dialogSlice = createSlice({
	name: 'dialog',
	initialState,
	reducers: {
		changeDialogIsOpen: (state, action: PayloadAction<IChangeDialogAction>) => {
			const { variant, isOpen, context } = action.payload
			state[variant] = { isOpen, context }
		},

		resetDialog: () => initialState,
	},
})

export const getDialogState = (variant: DialogVariants) => (state: RootState) => state.dialog[variant]

export const dialogPath = dialogSlice.name
export const dialogReducer = dialogSlice.reducer

export const { changeDialogIsOpen, resetDialog } = dialogSlice.actions
