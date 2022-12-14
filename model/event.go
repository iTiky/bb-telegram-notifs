package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

// EventType defines the Event type (e.g. "on PR open", "on new comment", etc.).
type EventType string

const (
	EventTypePROpen     EventType = "pr_open"
	EventTypePRApproved EventType = "pr_approved"
	EventTypePRRejected EventType = "pr_rejected"
	EventTypePRMerged   EventType = "pr_merged"
	EventTypePRUpdated  EventType = "pr_updated"
	EventTypeComment    EventType = "comment"
)

// Event defines an entry that is bridged to the Telegram bot.
// Hash field is used to prevent duplicate notifications.
type Event struct {
	bun.BaseModel `bun:"table:events,alias:e"`

	ID                int64      `bun:",pk,autoincrement"`
	Hash              string     // unique hash per event
	Type              EventType  // PR action type
	RecipientTgID     int64      // destination user Telegram ID
	RecipientTgChatID int64      // destination user Telegram chat ID
	SenderName        string     // source user name (Bitbucket's DisplayName)
	RepoProject       string     // source project
	RepoName          string     // source repo name
	PrID              int64      // source PR ID
	PrTitle           string     // source PR title
	PrURL             string     // source PR URL
	SendAck           bool       // whether the event was sent to the user (retry later if not)
	SendAt            *time.Time // event send time (nil if not sent yet)
	CreatedAt         time.Time  // event creation time (Bitbucket's Activity timestamp / PR creation timestamp)
}

// SetHash sets the event hash using the event data.
func (e *Event) SetHash() {
	hashStr := strings.Builder{}
	hashStr.WriteString(string(e.Type))
	hashStr.WriteString(strconv.FormatInt(e.RecipientTgID, 10))
	hashStr.WriteString(strconv.FormatInt(e.RecipientTgChatID, 10))
	hashStr.WriteString(e.SenderName)
	hashStr.WriteString(e.RepoProject)
	hashStr.WriteString(e.RepoName)
	hashStr.WriteString(strconv.FormatInt(e.PrID, 10))
	hashStr.WriteString(e.PrTitle)
	hashStr.WriteString(e.PrURL)
	hashStr.WriteString(e.CreatedAt.String())

	hash := sha1.Sum([]byte(hashStr.String()))

	e.Hash = hex.EncodeToString(hash[:])
}

// String implements the fmt.Stringer interface.
func (e Event) String() string {
	return fmt.Sprintf("%s/%s (%s): %s -> %d", e.RepoProject, e.RepoName, e.Type, e.SenderName, e.RecipientTgID)
}
