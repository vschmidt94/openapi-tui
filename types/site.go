package types

type Site struct {
	Name string
	Uri  string
	User string
}

func (s Site) FilterValue() string {
	return s.Name
}

func (s Site) Title() string {
	return s.Name
}

func (s Site) Description() string {
	return s.Uri
}
