package types

type Endpoint struct {
	Path        string
	Method      string
	Description string
}

func (ep Endpoint) FilterValue() string {
	return ep.Method + " " + ep.Path
}

func (ep Endpoint) GetTitle() string {
	return ep.Method + " " + ep.Path
}

func (ep Endpoint) GetDescription() string {
	return ep.Description
}
