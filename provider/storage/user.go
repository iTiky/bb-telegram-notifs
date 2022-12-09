package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/itiky/bb-telegram-notifs/model"
)

// CreateUser creates a new model.User.
func (st *Psql) CreateUser(ctx context.Context, tgID, tgChatID int64, tgLogin string) (*model.User, error) {
	if tgID == 0 {
		return nil, fmt.Errorf("telegram ID is empty")
	}
	if tgChatID == 0 {
		return nil, fmt.Errorf("telegram chatID is empty")
	}

	user := model.User{
		TgID:      tgID,
		TgLogin:   tgLogin,
		TgChatID:  tgChatID,
		Active:    false,
		CreatedAt: time.Now().UTC(),
	}
	_, err := st.db.NewInsert().Model(&user).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	return &user, nil
}

// GetUserByID returns the model.User for the given ID.
func (st *Psql) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := st.db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("selecting user by login: %w", err)
	}

	return &user, nil
}

// GetUserByTgID returns the model.User for the given Telegram ID.
func (st *Psql) GetUserByTgID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := st.db.NewSelect().
		Model(&user).
		Where("tg_id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("selecting user by login: %w", err)
	}

	return &user, nil
}

// GetUserByBBEmail returns the model.User for the given BitBucket email.
func (st *Psql) GetUserByBBEmail(ctx context.Context, bbEmail string) (*model.User, error) {
	var user model.User
	err := st.db.NewSelect().
		Model(&user).
		Where("bb_email = ?", bbEmail).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("selecting user by bbEmail: %w", err)
	}

	return &user, nil
}

// SetBBEmail sets the BitBucket email for the given user.
func (st *Psql) SetBBEmail(ctx context.Context, userID int64, bbEmail string) error {
	_, err := st.db.NewUpdate().
		Model(&model.User{}).
		Set("bb_email = ?", bbEmail).
		Where("id = ?", userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("setting bbEmail for user: %w", err)
	}

	return nil
}
