package models

import "github.com/goccy/go-json"

type Identity struct {
	Role   string           `json:"role"`
	Realm  string           `json:"realm"`
	UserId string           `json:"userId"`
	Roles  []*RoleWithRealm `json:"roles"`
}

func (r *Identity) String() (string, error) {
	data, err := json.Marshal(r)
	return string(data), err
}

func (r *Identity) Parse(data string) error {
	return json.Unmarshal([]byte(data), &r)
}
