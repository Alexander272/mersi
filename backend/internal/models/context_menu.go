package models

type ContextMenu struct {
	Id        string `json:"id" db:"id"`
	Position  int    `json:"position" db:"position"`
	SectionId string `json:"sectionId" db:"section_id"`
	Name      string `json:"name" db:"name"`
	Label     string `json:"label" db:"label"`
	Rule      string `json:"rule" db:"rule"`
}

type GetContextMenuDTO struct {
	SectionId string `json:"sectionId"`
	UserId    string `json:"userId"`
	Role      string `json:"role"`
}

type ContextMenuDTO struct {
	Id         string `json:"id" db:"id"`
	Position   int    `json:"position" db:"position" binding:"min=0"`
	SectionId  string `json:"sectionId" db:"section_id" binding:"required"`
	Name       string `json:"name" db:"name" binding:"required"`
	Label      string `json:"label" db:"label"`
	RuleItemId string `json:"ruleItemId" db:"rule_item_id" binding:"required"`
}

type DeleteContextMenuDTO struct {
	Id string `json:"id" db:"id"`
}

type CustomContextMenuDTO struct {
	Id          string `json:"id" db:"id"`
	UserId      string `json:"userId" db:"user_id"`
	ToolsMenuId string `json:"toolsMenuId" db:"tools_menu_id"`
}
