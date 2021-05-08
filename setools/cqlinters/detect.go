package cqlinters

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

// Detect finds all CQLinters that are used within the project.
func Detect(project api.Project) []api.CQLinter {
	res := []api.CQLinter{}
	for _, linter := range ByType {
		if DetectLinter(linter, project) {
			res = append(res, linter)
		}
	}
	return res
}

// DetectType returns true if this CQLinterType is being used in the project.
func DetectType(typ api.CQLinterType, project api.Project) bool {
	linter, ok := ByType[typ]
	return ok && DetectLinter(linter, project)
}

// DetectLinter returns true if this is CQLinter is being used in the project,
// i.e. when either linter.IsConfigured() is true, or when the project's main dependency manager has the linter in its dependencies.
func DetectLinter(linter api.CQLinter, project api.Project) bool {
	return linter.IsConfigured(project) ||
		(len(project.DepManagers) > 0 && project.DepManagers.Main().HasDependency(linter.DependencyName()))
}

//---------------------------------------------------------------------------------------

// FromConfig returns a list of CQLinters based on the CodeQualityConfig from the mllint configuration.
// Always returns a non-nil list of CQLinters containing all correctly named linters, even when there is an error.
func FromConfig(conf config.CodeQualityConfig) ([]api.CQLinter, error) {
	linters := []api.CQLinter{}
	notFound := []string{}

	for _, ltype := range conf.Linters {
		linter, ok := ByType[api.CQLinterType(strings.ToLower(ltype))]

		if ok {
			linters = append(linters, linter)
		} else {
			notFound = append(notFound, ltype)
		}
	}

	var err error
	if len(notFound) > 0 {
		err = fmt.Errorf("unknown code quality linters in mllint's config: %+v", notFound)
	}
	return linters, err
}
