package versioncontrol

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters/common"
)

func NewLinter() api.Linter {
	return common.NewCompositeLinter(categories.VersionControl.Name, &GitLinter{}, &DVCLinter{})
}
