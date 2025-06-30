package models

type ToolsMenu struct {
	Id            string `json:"id" db:"id"`
	Position      int    `json:"position" db:"position"`
	SectionId     string `json:"sectionId" db:"section_id"`
	Name          string `json:"name" db:"name"`
	Label         string `json:"label" db:"label"`
	Rule          string `json:"rule" db:"rule"`
	CanBeFavorite bool   `json:"canBeFavorite" db:"can_be_favorite"`
	Favorite      bool   `json:"favorite" db:"favorite"`
}

type GetToolsMenuDTO struct {
	SectionId string `json:"sectionId"`
	UserId    string `json:"userId"`
	Role      string `json:"role"`
}

type ChangeFavoriteDTO struct {
	Id       string `json:"id" db:"id"`
	Favorite bool   `json:"favorite" db:"favorite"`
	UserId   string
}

type ToolsMenuDTO struct {
	Id            string `json:"id" db:"id"`
	Position      int    `json:"position" db:"position"`
	SectionId     string `json:"sectionId" db:"section_id"`
	Name          string `json:"name" db:"name"`
	Label         string `json:"label" db:"label"`
	RuleItemId    string `json:"ruleItemId" db:"rule_item_id"`
	CanBeFavorite bool   `json:"canBeFavorite" db:"can_be_favorite"`
	Favorite      bool   `json:"favorite" db:"favorite"`
}

type DeleteToolsMenuDTO struct {
	Id string `json:"id" db:"id"`
}
