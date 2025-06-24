package pq_models

import "github.com/lib/pq"

type MenuDTO struct {
	Id          string         `json:"id" db:"id"`
	RoleId      string         `json:"-" db:"role_id"`
	RoleName    string         `json:"roleName" db:"name"`
	RoleLevel   int            `json:"roleLevel" db:"level"`
	RoleExtends pq.StringArray `json:"roleExtends" db:"extends"`
	RuleItemId  string         `json:"ruleItemId" db:"rule_item_id"`
}
