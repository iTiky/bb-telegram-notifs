package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/itiky/bb-telegram-notifs/model"
)

// CreateRepo creates a new model.Repo.
// Idempotent call.
func (st *Psql) CreateRepo(ctx context.Context, projectName, repoName string) error {
	if projectName == "" {
		return errors.New("project name is empty")
	}
	if repoName == "" {
		return errors.New("repo name is empty")
	}

	repo := model.Repo{
		Project:   projectName,
		Name:      repoName,
		CreatedAt: time.Now().UTC(),
	}

	_, err := st.db.NewInsert().Model(&repo).Exec(ctx)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			if pgErr.IntegrityViolation() {
				return nil
			}
		}
		return fmt.Errorf("inserting new repo: %w", err)
	}

	return nil
}

// GetRepo returns a model.Repo by project and name.
func (st *Psql) GetRepo(ctx context.Context, projectName, repoName string) (*model.Repo, error) {
	var repo model.Repo
	err := st.db.NewSelect().Model(&repo).Where("project = ? AND name = ?", projectName, repoName).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("selecting repo: %w", err)
	}

	return &repo, nil
}

// ListRepos returns a list of model.Repo.
func (st *Psql) ListRepos(ctx context.Context) ([]model.Repo, error) {
	var repos []model.Repo
	err := st.db.NewSelect().Model(&repos).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("selecting repos: %w", err)
	}

	return repos, nil
}
