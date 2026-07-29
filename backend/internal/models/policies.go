package models

type Access struct {
	Role   string   `json:"role"`
	Domain string   `json:"domain"`
	Perms  []string `json:"permissions"`
}

type PolicyResult struct {
	Perms []string `json:"permissions"`
}
