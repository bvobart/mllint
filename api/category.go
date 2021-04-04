package api

type Category string

func (c Category) Slug() string {
	return Slug(string(c))
}
