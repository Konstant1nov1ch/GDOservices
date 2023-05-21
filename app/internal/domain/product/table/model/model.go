package model

import "github.com/jackc/pgtype"

type Table struct {
	Id       int         `json:"id"`
	UserId   pgtype.UUID `json:"user_id"`
	Capacity int         `json:"capacity"`
}
