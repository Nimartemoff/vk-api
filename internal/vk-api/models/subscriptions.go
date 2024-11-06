package models

type Subscriptions struct {
	Users  []User  `json:"users"`
	Groups []Group `json:"groups"`
}
