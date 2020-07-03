package contribution

type CSVRow struct {
	CreatedAt    string
	Owner        string
	Repository   string
	Total        CSVColumn
	Contributors map[string]CSVColumn
}

type CSVColumn struct {
	PullRequest string
	Approve     string
	Comment     string
}

func (c CSVRow) Row(contributors []string) []string {
	var r []string
	r = append(r, c.CreatedAt)
	r = append(r, c.Owner)
	r = append(r, c.Repository)
	r = append(r, c.Total.Row()...)
	for _, name := range contributors {
		r = append(r, c.Contributors[name].Row()...)
	}
	return r
}

func (c CSVColumn) Row() []string {
	var r []string
	r = append(r, c.PullRequest)
	r = append(r, c.Approve)
	r = append(r, c.Comment)
	return r
}
