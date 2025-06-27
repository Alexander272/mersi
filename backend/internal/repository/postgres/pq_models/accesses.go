package pq_models

import "time"

type Accesses struct {
	Id        string    `db:"id"`
	RealmId   string    `db:"realm_id"`
	UserId    string    `db:"user_id"`
	SSOId     string    `db:"sso_id"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	RoleId    string    `db:"role_id"`
	RoleName  string    `db:"name"`
	Created   time.Time `db:"created_at"`
}
