package models

type Node struct {
	ID     int64    `json:"id"`
	Labels []string `json:"labels"`
}
