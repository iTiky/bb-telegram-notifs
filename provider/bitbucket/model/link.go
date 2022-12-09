package model

type Link struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type Links map[string][]Link
