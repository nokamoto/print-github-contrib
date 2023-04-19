package github

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	logger *zap.Logger
	sleep  time.Duration
}

func NewClient(logger *zap.Logger, sleep time.Duration) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{
		client: github.NewClient(tc),
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
	logger := c.logger.With(zap.String("org", org))
	logger.Debug("ListRepositoryByOrg")

	opts := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 30}}

	var repositories []*github.Repository
	for {
		repos, res, err := c.client.Repositories.ListByOrg(ctx, org, opts)

		time.Sleep(c.sleep)

		logger.Debug("ListByOrg", zap.Int("repos", len(repos)))
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, repos...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return repositories, nil
}

func (c *Client) ListPullRequest(ctx context.Context, owner string, repo string, start, end time.Time) ([]*github.PullRequest, error) {
	logger := c.logger.With(
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.Time("start", start),
		zap.Time("end", end),
	)
	logger.Debug("ListPullRequest")

	opts := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var pullRequests []*github.PullRequest
	for {
		prs, res, err := c.client.PullRequests.List(ctx, owner, repo, opts)

		time.Sleep(c.sleep)

		logger.Debug("List", zap.Int("prs", len(prs)))
		if err != nil {
			return nil, err
		}

		out := false
		for _, pr := range prs {
			prlog := logger.With(
				zap.Int("number", pr.GetNumber()),
				zap.Time("createdAt", pr.GetCreatedAt()),
				zap.String("title", pr.GetTitle()),
			)
			if pr.GetCreatedAt().Before(start) {
				prlog.Debug("[out-of-range] PR createdAt before start")
				out = true
				continue
			}
			if pr.GetCreatedAt().Before(end) {
				prlog.Debug("PR createdAt before end")
				pullRequests = append(pullRequests, pr)
				continue
			}
			prlog.Debug("[out-of-range] PR createdAt after end")
		}

		if out {
			break
		}

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return pullRequests, nil
}
