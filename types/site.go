package types

type Site struct {
	Name         string
	Uri          string
	User         string
	RequiresAuth bool
}

func (s Site) FilterValue() string {
	return s.Name
}

func (s Site) GetTitle() string {
	return s.Name
}

func (s Site) GetDescription() string {
	return s.Uri
}

func (s Site) GetUser() string {
	return s.User
}

func (s Site) GetRequiresAuth() bool {
	return s.RequiresAuth
}
