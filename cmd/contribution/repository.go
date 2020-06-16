package contribution

type Repository struct {
	Name         string
	Contributors map[string]*Contribution
}
