package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/itiky/bb-telegram-notifs/model"
)

// GetKeyValue returns a model.KeyValue by key.
func (st *Psql) GetKeyValue(ctx context.Context, key string) (*string, error) {
	kv := model.KeyValue{
		ID: key,
	}
	err := st.db.NewSelect().
		Model(&kv).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("selecting key-value: %w", err)
	}

	return &kv.Value, nil
}

// SetKeyValue sets a model.KeyValue (upsert).
func (st *Psql) SetKeyValue(ctx context.Context, key, value string) error {
	kv := model.KeyValue{
		ID:        key,
		Value:     value,
		UpdatedAt: time.Now().UTC(),
	}
	_, err := st.db.NewInsert().
		Model(&kv).
		On("CONFLICT (id) DO UPDATE").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("inserting key-value: %w", err)
	}

	return nil
}
