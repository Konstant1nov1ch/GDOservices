package model

type Table struct {
	Id          string `json:"id         "`
	TableId     string `json:"table_id   "`
	CategoryId  string `json:"category_id"`
	Deadline    string `json:"deadline   "`
	Title       string `json:"title      "`
	Description string `json:"description"`
}
