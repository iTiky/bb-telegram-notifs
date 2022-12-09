package model

import (
	"time"
)

type PRActivity struct {
	ID          int64 `json:"id"`
	CreatedDate int64 `json:"createdDate"`
	User        struct {
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
	} `json:"user"`
	Action string `json:"action"`
}

func (a PRActivity) CreatedAt() time.Time {
	return TimestampToTime(a.CreatedDate)
}
