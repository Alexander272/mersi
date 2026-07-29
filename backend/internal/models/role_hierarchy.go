package models

type RoleHierarchy struct {
	Role       Role `json:"childRole"`
	ParentRole Role `json:"parentRole"`
}

type RoleHierarchyDTO struct {
	ParentRoleId string `json:"parentRoleId" db:"parent_role_id"`
	RoleId       string `json:"childRoleId" db:"child_role_id"`
}

type GetRolesInheritance struct {
	Roles []string
}

type SyncRoleInheritance struct {
	Role       string
	ParentRole string
	Realm      string
}
