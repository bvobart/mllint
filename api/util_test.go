package api_test

import (
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"

	"github.com/stretchr/testify/require"
)

func TestSlugIsLowercaseDashed(t *testing.T) {
	for _, cat := range categories.All {
		require.Regexp(t, "^([a-z]+-)*[a-z]+$", api.Slug(cat.Name))
	}
}
