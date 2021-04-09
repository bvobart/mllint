package api

type Category struct {
	Name        string
	Slug        string
	Description string
}

func (c Category) String() string {
	return c.Name
}
