package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/pkg/logging"
)

//go:embed migration/*.sql
var dbMigrations embed.FS

// Psql is a PostgreSQL storage provider.
type Psql struct {
	db *bun.DB
}

// NewPsql creates a new Psql instance.
func NewPsql(ctx context.Context) (*Psql, error) {
	// Connection
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBName),
		viper.GetString(config.DBSSLMode),
	)
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbDSN)))

	st := &Psql{
		db: bun.NewDB(pgDB, pgdialect.New()),
	}
	if err := st.db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	// Migrations
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(dbMigrations); err != nil {
		return nil, fmt.Errorf("migration: discover: %w", err)
	}

	migrator := migrate.NewMigrator(st.db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return nil, fmt.Errorf("migration: init: %w", err)
	}

	if err := migrator.Lock(ctx); err != nil {
		return nil, fmt.Errorf("migration: lock: %w", err)
	}
	defer migrator.Unlock(ctx)

	migrationGroup, err := migrator.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("migration: migrate: %w", err)
	}
	st.Logger(ctx).Debug().Msgf("Migrations applied: %s", migrationGroup.String())

	return st, nil
}

// Logger returns a logger instance.
func (st *Psql) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.KeyProvider, "psql").Logger()

	return &logger
}
