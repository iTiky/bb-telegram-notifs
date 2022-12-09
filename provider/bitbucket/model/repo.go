package model

type Repo struct {
	ID          int64   `json:"id"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Project     Project `json:"project"`
	Public      bool    `json:"public"`
	Links       Links   `json:"links"`
}
