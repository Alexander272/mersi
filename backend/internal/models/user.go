package models

type User struct {
	ID   string   `json:"id" db:"id"`
	Name string   `json:"name" db:"name"`
	Role string   `json:"role"`
	Menu []string `json:"menu"`
	// Filters []*SIFilter `json:"filters"`

	Roles []*RoleWithRealm `json:"-"`

	AccessToken  string `json:"token"`
	RefreshToken string `json:"-"`
}

type UserData struct {
	ID        string `json:"id" db:"id"`
	SSOID     string `json:"ssoId" db:"sso_id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email" db:"email"`
}

type GetByRealmDTO struct {
	RealmID string `json:"realmId" binding:"required"`
	Include bool   `json:"include"`
}

type GetByAccessDTO struct {
	RealmID string `json:"realmId" binding:"required"`
	Role    string `json:"role"`
}

// type KeycloakUser struct {
// 	Id        string `json:"id"`
// 	Username  string `json:"username"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	Email     string `json:"email"`
// }
