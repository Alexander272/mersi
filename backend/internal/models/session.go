package models

type SignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Realm    string `json:"realm"`
}

type RefreshDTO struct {
	Token string `json:"token" binding:"required"`
	Realm string `json:"realm"`
}
