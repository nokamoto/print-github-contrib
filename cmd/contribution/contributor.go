package contribution

import (
	"fmt"
	"github.com/google/go-github/v32/github"
)

type Contributor struct {
	Name        string
	CreatedAt   string
	PullRequest int
	Approve     int
	Comment     int
}

type ContributorKey struct {
	Name      string
	CreatedAt string
}

func (k ContributorKey) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%s", k.Name, k.CreatedAt)), nil
}

type Contributors map[ContributorKey]*Contributor

func (c Contributors) Get(name string, pr *github.PullRequest) *Contributor {
	month := pr.GetCreatedAt().Format("2006-01")
	k := ContributorKey{
		Name:      name,
		CreatedAt: month,
	}
	_, ok := c[k]
	if !ok {
		c[k] = &Contributor{
			Name:      name,
			CreatedAt: month,
		}
	}
	return c[k]
}
