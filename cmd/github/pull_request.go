package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"go.uber.org/zap"
	"time"
)

type Contribution struct {
	PullRequest *github.PullRequest
	Reviews     []*github.PullRequestReview
	Comments    []*github.PullRequestComment
}

func (c *Client) ListReview(ctx context.Context, owner string, repo string, pr *github.PullRequest) ([]*github.PullRequestReview, error) {
	logger := c.logger.With(zap.String("owner", owner), zap.String("repo", repo), zap.Int("pr", pr.GetNumber()))
	logger.Debug("ListReview")

	opts := &github.ListOptions{PerPage: 30}

	var reviews []*github.PullRequestReview
	for {
		rs, res, err := c.client.PullRequests.ListReviews(ctx, owner, repo, pr.GetNumber(), opts)

		time.Sleep(c.sleep)

		logger.Debug("ListReviews", zap.Int("reviews", len(rs)))
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, rs...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return reviews, nil
}

func (c *Client) ListComment(ctx context.Context, owner string, repo string, pr *github.PullRequest) ([]*github.PullRequestComment, error) {
	logger := c.logger.With(zap.String("owner", owner), zap.String("repo", repo), zap.Int("pr", pr.GetNumber()))
	logger.Debug("ListComment")

	opts := &github.PullRequestListCommentsOptions{ListOptions: github.ListOptions{PerPage: 30}}

	var comments []*github.PullRequestComment
	for {
		cs, res, err := c.client.PullRequests.ListComments(ctx, owner, repo, pr.GetNumber(), opts)

		time.Sleep(c.sleep)

		logger.Debug("ListComments", zap.Int("comments", len(cs)))
		if err != nil {
			return nil, err
		}

		comments = append(comments, cs...)

		if res.NextPage == 0 {
			break
		}

		opts.Page = res.NextPage
	}

	return comments, nil
}

func (c *Client) ListContribution(ctx context.Context, owner string, repo string, prs []*github.PullRequest) ([]*Contribution, error) {
	var contrib []*Contribution

	for _, pr := range prs {
		reviews, err := c.ListReview(ctx, owner, repo, pr)
		if err != nil {
			return nil, err
		}

		comments, err := c.ListComment(ctx, owner, repo, pr)
		if err != nil {
			return nil, err
		}

		contrib = append(contrib, &Contribution{
			PullRequest: pr,
			Reviews:     reviews,
			Comments:    comments,
		})
	}

	return contrib, nil
}
