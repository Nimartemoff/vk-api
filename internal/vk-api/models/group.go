package models

type Group struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
