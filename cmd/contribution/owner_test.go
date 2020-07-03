package contribution

import (
	"fmt"
	gogithub "github.com/google/go-github/v32/github"
	"github.com/nokamoto/print-github-contrib/cmd/github"
	"time"
)

func ExampleOwner_CSV() {
	o := NewOwner("foo", time.Now(), time.Now())

	name := func(s string) *string {
		return &s
	}

	bob := gogithub.User{
		Login: name("bob"),
	}

	alice := gogithub.User{
		Login: name("alice"),
	}

	fred := gogithub.User{
		Login: name("fred"),
	}

	barney := gogithub.User{
		Login: name("barney"),
	}

	approved := "APPROVED"

	utc, _ := time.LoadLocation("UTC")
	now := time.Date(2020, 6, 17, 0, 0, 0, 0, utc)

	lastMonth := time.Date(2020, 5, 17, 0, 0, 0, 0, utc)

	bar := NewRepository("bar", []*github.Contribution{
		{
			PullRequest: &gogithub.PullRequest{
				CreatedAt: &now,
				User:      &bob,
			},
			Reviews: []*gogithub.PullRequestReview{
				{
					User:  &alice,
					State: &approved,
				},
			},
			Comments: []*gogithub.PullRequestComment{
				{
					User: &fred,
				},
				{
					User: &barney,
				},
			},
			IssueComments: []*gogithub.IssueComment{
				{
					User: &fred,
				},
			},
		},
		{
			PullRequest: &gogithub.PullRequest{
				CreatedAt: &now,
				User:      &barney,
			},
			Reviews: []*gogithub.PullRequestReview{
				{
					User:  &bob,
					State: &approved,
				},
			},
			Comments: []*gogithub.PullRequestComment{
				{
					User: &fred,
				},
				{
					User: &barney,
				},
			},
		},
		{
			PullRequest: &gogithub.PullRequest{
				CreatedAt: &lastMonth,
				User:      &alice,
			},
			Reviews: []*gogithub.PullRequestReview{
				{
					User:  &bob,
					State: &approved,
				},
			},
		},
	})

	o.AddRepository(bar)

	csv, err := o.CSV()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(csv)

	// Output:
	// created_at,owner,repository,pull_request,approve,comment,alice.pull_request,alice.approve,alice.comment,barney.pull_request,barney.approve,barney.comment,bob.pull_request,bob.approve,bob.comment,fred.pull_request,fred.approve,fred.comment
	// 2020-05,foo,bar,1,1,0,1,0,0,,,,0,1,0,,,
	// 2020-06,foo,bar,2,2,5,0,1,0,1,0,2,1,1,0,0,0,3
}
