package dependencymgmt

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/setools/depmanagers"
)

func NewLinter() api.Linter {
	return &DependenciesLinter{}
}

// This linter relates to the best practice of using proper dependency management,
// as found to be a major obstacle towards reproducibility of ML projects in https://arxiv.org/abs/2103.04146
type DependenciesLinter struct{}

func (l *DependenciesLinter) Name() string {
	return categories.DependencyMgmt.Name
}

func (l *DependenciesLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleUse, &RuleSingle} // TODO: add the rest
}

func (l *DependenciesLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	managers := project.DepManagers

	if len(managers) == 0 {
		report.Scores[RuleUse] = 0
		report.Scores[RuleSingle] = 0
		return report, nil
	}

	types := types(managers)
	switch {
	case contains(types, depmanagers.TypePipenv) || contains(types, depmanagers.TypePoetry):
		report.Scores[RuleUse] = 100
	case contains(types, depmanagers.TypeRequirementsTxt):
		report.Scores[RuleUse] = 20 // it's better than nothing, but not recommended
		report.Details[RuleUse] = DetailsNoRequirementsTxt
	case contains(types, depmanagers.TypeSetupPy):
		report.Scores[RuleUse] = 30 // it's better than nothing and slightly better than a requirements.txt, but still not recommended.
		report.Details[RuleUse] = DetailsNoSetupPy
	default:
		report.Scores[RuleUse] = 0
		report.Details[RuleUse] = fmt.Sprintf("Your project is somehow using a dependency manager that mllint recognises, but cannot score: %s.\n\nPlease create an issue on mllint's GitHub :)", types)
	}

	if len(managers) == 1 {
		report.Scores[RuleSingle] = 100
		return report, nil
	}

	report.Scores[RuleSingle] = 0

	details := strings.Builder{}
	details.WriteString(fmt.Sprintf("Your project was found to be using dependency managers: %+v\n\n", types))
	switch {
	case containsAll(types, depmanagers.TypeRequirementsTxt, depmanagers.TypeSetupPy):
		details.WriteString(DetailsRequirementsTxtSetupPy)
	case containsAll(types, depmanagers.TypeRequirementsTxt, depmanagers.TypePipenv):
		details.WriteString(DetailsRequirementsTxtPipenv)
	case containsAll(types, depmanagers.TypeRequirementsTxt, depmanagers.TypePoetry):
		details.WriteString(DetailsRequirementsTxtPoetry)
	case containsAll(types, depmanagers.TypePipenv, depmanagers.TypeSetupPy):
		details.WriteString(DetailsPipenvSetupPy)
	case containsAll(types, depmanagers.TypePoetry, depmanagers.TypeSetupPy):
		details.WriteString(DetailsPoetrySetupPy)
	default:
		details.WriteString("Pick the one most suited for you, your project and your team, then stick with it.")
	}
	report.Details[RuleSingle] = details.String()

	return report, nil
}

func types(managers []api.DependencyManager) []api.DependencyManagerType {
	types := []api.DependencyManagerType{}
	for _, manager := range managers {
		types = append(types, manager.Type())
	}
	return types
}

func contains(types []api.DependencyManagerType, target api.DependencyManagerType) bool {
	for _, typ := range types {
		if typ == target {
			return true
		}
	}
	return false
}

func containsAll(types []api.DependencyManagerType, targets ...api.DependencyManagerType) bool {
	for _, target := range targets {
		if !contains(types, target) {
			return false
		}
	}
	return true
}
