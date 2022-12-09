package model

import (
	"time"
)

const (
	PRStateOpen = "OPEN"
)

type PR struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	State       string   `json:"state"`
	CreatedDate int64    `json:"createdDate"`
	Author      PRUser   `json:"author"`
	Reviewers   []PRUser `json:"reviewers"`
	Links       Links    `json:"links"`
}

func (pr PR) CreatedAt() time.Time {
	return TimestampToTime(pr.CreatedDate)
}

func (pr PR) SelfLink() *Link {
	l, found := pr.Links["self"]
	if !found {
		return nil
	}

	if len(l) > 0 {
		return &l[0]
	}

	return nil
}

type PRUser struct {
	User struct {
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
	} `json:"user"`
	Approved bool   `json:"approved"`
	Status   string `json:"status"`
}
