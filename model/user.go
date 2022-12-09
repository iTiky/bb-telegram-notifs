package model

import (
	"time"

	"github.com/uptrace/bun"
)

// User defines a model to track Telegram and BitBucket user specifics.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64  `bun:",pk,autoincrement"`
	TgID      int64  // Telegram user ID
	TgLogin   string // Telegram user name (could be empty)
	TgChatID  int64  // Telegram chat ID (equals to TgID if user is not in a group)
	BbEmail   string // BitBucket user email
	Active    bool   // if not active, user is not authorized to use the bot
	CreatedAt time.Time
}
