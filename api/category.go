package api

import (
	"strings"
)

type Category string

func (c Category) Slug() string {
	return strings.ReplaceAll(strings.ToLower(string(c)), " ", "-")
}
