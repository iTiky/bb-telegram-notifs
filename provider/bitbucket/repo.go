package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itiky/bb-telegram-notifs/provider/bitbucket/model"
)

// ListRepos lists all repos for the current project.
func (c *Client) ListRepos(ctx context.Context) ([]model.Repo, error) {
	const endpointFmt = "projects/%s/repos" // {1} - project key

	endpoint := fmt.Sprintf(endpointFmt, c.bbProject)

	var list []model.Repo
	valuesHandler := func(v json.RawMessage) error {
		var repos []model.Repo
		if err := json.Unmarshal(v, &repos); err != nil {
			return err
		}
		list = append(list, repos...)
		return nil
	}

	if err := c.doPageAllRequest(ctx, endpoint, valuesHandler); err != nil {
		return nil, fmt.Errorf("list repos: %w", err)
	}

	return list, nil
}
