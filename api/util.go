package api

import "strings"

// Slug returns the given string, lowercased and with all spaces replaced by dashes.
func Slug(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}
