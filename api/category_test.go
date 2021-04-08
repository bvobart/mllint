package api_test

import (
	"testing"

	"github.com/bvobart/mllint/api/categories"
	"github.com/stretchr/testify/require"
)

func TestSlugIsLowercaseDashed(t *testing.T) {
	for _, cat := range categories.All {
		require.Regexp(t, "^([a-z]+-)*[a-z]+$", cat.Slug())
	}
}
