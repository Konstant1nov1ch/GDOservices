package model

type Category struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	BackColor string `json:"back_color"`
	WordColor string `json:"word_color"`
}
