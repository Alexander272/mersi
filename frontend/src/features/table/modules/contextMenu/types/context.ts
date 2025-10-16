export interface IContextMenu {
	id: string
	position: number
	sectionId: string
	name: string
	label: string
	rule: string
}

export interface IContextMenuDTO {
	id: string
	sectionId: string
	position: number
	name: string
	label: string
	rule: string
	ruleItemId: string
}
