package versioncontrol

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/categories"
	"github.com/bvobart/mllint/linters/common"
)

func NewLinter() api.Linter {
	return common.NewCompositeLinter(string(categories.VersionControl), &GitLinter{}, &DVCLinter{})
}
