package model

import (
	"time"

	"github.com/uptrace/bun"
)

// Repo defines a repository model.
// Used to link Subscriptions and Users.
type Repo struct {
	bun.BaseModel `bun:"table:repos,alias:r"`

	ID        int64  `bun:",pk,autoincrement"`
	Project   string // BitBucket's project name
	Name      string // BitBucket's repository name
	CreatedAt time.Time
}

// String implements the fmt.Stringer interface.
func (r Repo) String() string {
	return r.Project + "/" + r.Name
}
