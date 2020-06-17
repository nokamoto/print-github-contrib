package contribution

import (
	"fmt"
	"github.com/nokamoto/print-github-contrib/cmd/github"
	"sort"
)

type Repository struct {
	Name         string
	PullRequest  int
	Contributors Contributors
}

func NewRepository(name string, contrib []*github.Contribution) *Repository {
	repo := &Repository{
		Name:         name,
		Contributors: make(Contributors),
	}

	for _, c := range contrib {
		repo.PullRequest += 1

		repo.Contributors.Get(c.PullRequest.GetUser().GetLogin(), c.PullRequest).PullRequest += 1

		approve := make(map[string]struct{})
		for _, review := range c.Reviews {
			name := review.GetUser().GetLogin()

			if review.GetState() == "APPROVED" {
				approve[name] = struct{}{}
			}

			repo.Contributors.Get(name, c.PullRequest).Reviews += 1
		}

		for name := range approve {
			repo.Contributors.Get(name, c.PullRequest).Approve += 1
		}

		for _, comment := range c.Comments {
			repo.Contributors.Get(comment.GetUser().GetLogin(), c.PullRequest).Comment += 1
		}
	}

	return repo
}

func (r *Repository) Rows(owner string, contributorsName []string) []CSVRow {
	createdAt := make(map[string]struct{})
	for _, c := range r.Contributors {
		createdAt[c.CreatedAt] = struct{}{}
	}

	var sortedCreatedAt []string
	for at := range createdAt {
		sortedCreatedAt = append(sortedCreatedAt, at)
	}
	sort.Strings(sortedCreatedAt)

	var rows []CSVRow
	for _, at := range sortedCreatedAt {
		row := CSVRow{
			CreatedAt:    at,
			Owner:        owner,
			Repository:   r.Name,
			Total:        CSVColumn{},
			Contributors: map[string]CSVColumn{},
		}

		totalPullRequest := 0
		totalApprove := 0
		totalReviews := 0
		totalComment := 0

		for _, c := range r.Contributors {
			if c.CreatedAt == at {
				totalPullRequest += c.PullRequest
				totalApprove += c.Approve
				totalReviews += c.Reviews
				totalComment += c.Comment

				row.Contributors[c.Name] = CSVColumn{
					PullRequest: fmt.Sprintf("%d", c.PullRequest),
					Approve:     fmt.Sprintf("%d", c.Approve),
					Review:      fmt.Sprintf("%d", c.Reviews),
					Comment:     fmt.Sprintf("%d", c.Comment),
				}
			}
		}

		row.Total = CSVColumn{
			PullRequest: fmt.Sprintf("%d", totalPullRequest),
			Approve:     fmt.Sprintf("%d", totalApprove),
			Review:      fmt.Sprintf("%d", totalReviews),
			Comment:     fmt.Sprintf("%d", totalComment),
		}

		rows = append(rows, row)
	}

	return rows
}
