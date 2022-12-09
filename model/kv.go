package model

import (
	"time"

	"github.com/uptrace/bun"
)

// KeyValue defines the model to store key-value pairs.
type KeyValue struct {
	bun.BaseModel `bun:"table:kvs,alias:kv"`

	ID        string `bun:",pk"`
	Value     string
	UpdatedAt time.Time
}
