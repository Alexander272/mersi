export interface IToolsMenu {
	id: string
	position: number
	sectionId: string
	name: string
	label: string
	rule: string
	canBeFavorite: boolean
	favorite: boolean
}

export interface IToggleToolsMenu {
	id: string
	favorite: boolean
}
