package model

type Table struct {
	Id       int    `json:"id"`
	UserID   string `json:"user_id"`
	Capacity int    `json:"capacity"`
}
