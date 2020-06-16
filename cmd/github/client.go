package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	client *github.Client
	logger *zap.Logger
	sleep  time.Duration
}

func NewClient(logger *zap.Logger, sleep time.Duration) *Client {
	return &Client{
		client: github.NewClient(nil),
		logger: logger,
		sleep:  sleep,
	}
}

func NewEnterpriseClient(logger *zap.Logger, baseURL string, sleep time.Duration) (*Client, error) {
	client, err := github.NewEnterpriseClient(baseURL, "", nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		logger: logger,
		sleep:  sleep,
	}, nil
}

func (c *Client) ListRepositoryByOrg(ctx context.Context, org string) ([]*github.Repository, error) {
	c.logger.Debug("ListRepositoryByOrg", zap.String("org", org))

	opts := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 30}}

	var repositories []*github.Repository
	for {
		repos, res, err := c.client.Repositories.ListByOrg(ctx, org, opts)
		c.logger.Debug("ListByOrg", zap.Int("repos", len(repos)))
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repos...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage

		time.Sleep(c.sleep)
	}

	return repositories, nil
}
