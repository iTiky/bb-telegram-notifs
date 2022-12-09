package bitbucket

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/itiky/bb-telegram-notifs/pkg/config"
	"github.com/itiky/bb-telegram-notifs/provider/bitbucket/model"
)

// Client is a BitBucket API client.
// Client only works with a single BitBucket project.
type Client struct {
	token      string
	hostURL    *url.URL
	bbProject  string
	httpClient *http.Client
}

// NewClient creates a new BitBucket API client.
// Context is used to limit the ping duration.
func NewClient(ctx context.Context) (*Client, error) {
	token := viper.GetString(config.BBToken)

	hostBz := viper.GetString(config.BBHost)
	host, err := url.Parse(hostBz)
	if err != nil {
		return nil, fmt.Errorf("parse host (%s): %w", hostBz, err)
	}

	bbProject := viper.GetString(config.BBProject)
	if bbProject == "" {
		return nil, fmt.Errorf("bitBucket project is not defined")
	}

	c := &Client{
		token:     token,
		hostURL:   host,
		bbProject: bbProject,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 50,
				MaxIdleConns:        50,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	if err := c.Ping(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

// Ping pings the BitBucket API.
func (c *Client) Ping(ctx context.Context) error {
	const endpoint = "application-properties"

	req := c.buildRequest(ctx, endpoint)
	if err := c.doRequest(req, nil); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

// doPageAllRequest performs a set of request to read out all the pages fot the specified endpoint.
func (c *Client) doPageAllRequest(ctx context.Context, endpoint string, handleValues func(v json.RawMessage) error, queryKeyVals ...string) error {
	type pageData struct {
		model.Page
		Values json.RawMessage `json:"values"`
	}

	offset, limit, hasNextPage := 0, 50, true
	for hasNextPage {
		f := func() error {
			queryParams := append(
				queryKeyVals,
				"start", strconv.Itoa(offset),
				"limit", strconv.Itoa(limit),
			)
			req := c.buildRequest(ctx, endpoint, queryParams...)

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("request (%s): %w", req.URL.String(), err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("request (%s): reading response body: %w", req.URL.String(), err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("request (%s): unexpected status code: %d", req.URL.String(), resp.StatusCode)
			}

			var page pageData
			if err := json.Unmarshal(body, &page); err != nil {
				return fmt.Errorf("request (%s): page response unmarshal: %w", req.URL.String(), err)
			}

			if err := handleValues(page.Values); err != nil {
				return fmt.Errorf("request (%s): values response unmarshal: %w", req.URL.String(), err)
			}

			if !page.IsLastPage {
				offset = page.NextPageStart
			}
			hasNextPage = false

			return nil
		}
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

// doRequest performs a request to the specified endpoint unmarshalling the response into the specified data pointer.
func (c *Client) doRequest(req *http.Request, bodyObjPtr interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request (%s): %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("request (%s): reading response body: %w", req.URL.String(), err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request (%s): unexpected status code: %d", req.URL.String(), resp.StatusCode)
	}

	if bodyObjPtr == nil {
		return nil
	}

	if err := json.Unmarshal(body, bodyObjPtr); err != nil {
		return fmt.Errorf("request (%s): unmarshalling response body to %T: %w", req.URL.String(), bodyObjPtr, err)
	}

	return nil
}

// buildRequest builds a new request to the specified endpoint with the specified query parameters.
func (c *Client) buildRequest(ctx context.Context, endpoint string, queryKeyVals ...string) *http.Request {
	reqURL := *c.hostURL
	reqURL.Path = path.Join(reqURL.Path, endpoint)

	reqQuery := reqURL.Query()
	for i := 0; i < len(queryKeyVals)-1; i += 2 {
		reqQuery.Set(queryKeyVals[i], queryKeyVals[i+1])
	}
	reqURL.RawQuery = reqQuery.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	return req
}
