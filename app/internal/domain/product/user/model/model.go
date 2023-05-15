package model

type User struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Pwd           string `json:"pwd"`
	PaymentStatus string `json:"payment_status"`
}
