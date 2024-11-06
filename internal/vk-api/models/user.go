package models

type User struct {
	ID            uint64 `json:"id"`
	ScreenName    string `json:"screen_name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Sex           byte   `json:"sex"`
	City          City   `json:"city"`
	Followers     []User
	Subscriptions Subscriptions
}

type City struct {
	Title string `json:"title"`
}
