package model

import (
	"time"
)

func TimestampToTime(ts int64) time.Time {
	tsSec := ts / 1000
	tsNano := (ts % 1000) * 1000000

	return time.Unix(tsSec, tsNano).UTC()
}
