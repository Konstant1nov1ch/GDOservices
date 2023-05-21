package model

import "github.com/jackc/pgtype"

type User struct {
	Id            pgtype.UUID `json:"id"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	Pwd           string      `json:"pwd"`
	PaymentStatus bool        `json:"payment_status"`
}
