package flags

import "fmt"

type Output struct {
	Value string
}

const (
	JSON = "json"
	CSV  = "csv"
)

func (o *Output) String() string {
	return o.Value
}

func (o *Output) Set(s string) error {
	switch s {
	case JSON:
		o.Value = JSON
	case CSV:
		o.Value = CSV
	default:
		return fmt.Errorf("%s is not %s or %s", s, JSON, CSV)
	}
	return nil
}

func (*Output) Type() string {
	return "string"
}
