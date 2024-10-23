package models

type Subscription struct {
	Users  []User  `json:"users"`
	Groups []Group `json:"groups"`
}
