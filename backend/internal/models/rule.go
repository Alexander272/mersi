package models

type Rule struct {
	Id       string `json:"id" db:"id"`
	RoleId   string `json:"-" db:"role_id"`
	RoleName string `json:"roleName" db:"name"`
	// RoleLevel   int      `json:"roleLevel" db:"level"`
	// RoleExtends []string `json:"roleExtends" db:"extends"`
	ItemId     string `json:"ruleItemId" db:"rule_item_id"`
	ItemName   string `json:"ruleItemName" db:"item_name"`
	ItemMethod string `json:"ruleItemMethod" db:"method"`
}

type RuleFull struct {
	Id string `json:"id" db:"id"`
	// RoleId    string     `json:"-" db:"role_id"`
	// Role      string     `json:"role" db:"role"`
	Role      *RoleFull   `json:"role"`
	RuleItems []*RuleItem `json:"ruleItems"`
}

type RuleDTO struct {
	Id         string `json:"id"`
	RoleId     string `json:"roleId"`
	RuleItemId string `json:"ruleItemId"`
}
