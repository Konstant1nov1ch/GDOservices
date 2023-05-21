package model

import "github.com/jackc/pgtype"

type Note struct {
	Id          int         `json:"id"`
	TableId     int         `json:"table_id"`
	CategoryId  int         `json:"category_id"`
	Deadline    pgtype.Time `json:"deadline"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}
