package contribution

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Owner struct {
	Name         string
	Start        time.Time
	End          time.Time
	Repositories []*Repository
	Contributors map[string]struct{}
}

func NewOwner(name string, start, end time.Time) *Owner {
	return &Owner{Name: name, Start: start, End: end, Contributors: map[string]struct{}{}}
}

func (o *Owner) AddRepository(repo *Repository) {
	o.Repositories = append(o.Repositories, repo)

	for _, c := range repo.Contributors {
		o.Contributors[c.Name] = struct{}{}
	}
}

func (o *Owner) JSON() (string, error) {
	bs, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(bs), err
}

func (o *Owner) CSV() (string, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	var cs []string
	for c := range o.Contributors {
		cs = append(cs, c)
	}
	sort.Strings(cs)

	header := CSVRow{
		CreatedAt:  "created_at",
		Owner:      "owner",
		Repository: "repository",
		Total: CSVColumn{
			PullRequest: "pull_request",
			Approve:     "approve",
			Review:      "review",
			Comment:     "comment",
		},
		Contributors: map[string]CSVColumn{},
	}

	for _, c := range cs {
		header.Contributors[c] = CSVColumn{
			PullRequest: fmt.Sprintf("%s.pull_request", c),
			Approve:     fmt.Sprintf("%s.approve", c),
			Review:      fmt.Sprintf("%s.review", c),
			Comment:     fmt.Sprintf("%s.comment", c),
		}
	}

	err := w.Write(header.Row(cs))
	if err != nil {
		return "", err
	}

	for _, repo := range o.Repositories {
		for _, row := range repo.Rows(o.Name, cs) {
			err = w.Write(row.Row(cs))
			if err != nil {
				return "", err
			}
		}
	}

	w.Flush()

	return buf.String(), nil
}
