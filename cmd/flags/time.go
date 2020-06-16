package flags

import "time"

type Time struct {
	time.Time
}

const (
	layout = "2006-01-02"
)

func (t *Time) String() string {
	return t.Format(layout)
}

func (t *Time) Set(s string) error {
	v, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

func (t *Time) Type() string {
	return "string"
}
