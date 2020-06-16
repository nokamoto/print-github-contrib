package contribution

import "time"

type Owner struct {
	Name         string
	Start        time.Time
	End          time.Time
	Repositories map[string]*Repository
}
