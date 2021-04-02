package linters

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/categories"
	"github.com/bvobart/mllint/linters/versioncontrol"
)

var ByCategory = map[api.Category]api.Linter{
	categories.VersionControl: versioncontrol.NewLinter(),
}
