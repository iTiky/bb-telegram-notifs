package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itiky/bb-telegram-notifs/provider/bitbucket/model"
)

// ListPRsForRepo lists all PRs for the given repo.
func (c *Client) ListPRsForRepo(ctx context.Context, repoName, prState string) ([]model.PR, error) {
	const endpointFmt = "projects/%s/repos/%s/pull-requests" // {1} - project key, {2} - repo slug

	endpoint := fmt.Sprintf(endpointFmt, c.bbProject, repoName)

	var list []model.PR
	valuesHandler := func(v json.RawMessage) error {
		var prs []model.PR
		if err := json.Unmarshal(v, &prs); err != nil {
			return err
		}
		list = append(list, prs...)
		return nil
	}

	if err := c.doPageAllRequest(ctx, endpoint, valuesHandler, "state", prState); err != nil {
		return nil, fmt.Errorf("list PRs: %w", err)
	}

	return list, nil
}

// ListPRActivity lists all PR activity for the given repo.
func (c *Client) ListPRActivity(ctx context.Context, repoName string, repoBbID int64) ([]model.PRActivity, error) {
	const endpointFmt = "projects/%s/repos/%s/pull-requests/%d/activities" // {1} - project key, {2} - repo slug, {3} - PR ID

	endpoint := fmt.Sprintf(endpointFmt, c.bbProject, repoName, repoBbID)

	var list []model.PRActivity
	valuesHandler := func(v json.RawMessage) error {
		var activities []model.PRActivity
		if err := json.Unmarshal(v, &activities); err != nil {
			return err
		}
		list = append(list, activities...)
		return nil
	}

	if err := c.doPageAllRequest(ctx, endpoint, valuesHandler); err != nil {
		return nil, fmt.Errorf("list PRs: %w", err)
	}

	return list, nil
}
