package dependencymgmt

import (
	"fmt"
	"math"
	"strings"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils/markdowngen"
	"github.com/juliangruber/go-intersect"
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
	return []*api.Rule{&RuleUse, &RuleSingle, &RuleUseDev}
}

func (l *DependenciesLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	managers := project.DepManagers

	l.ScoreRuleUse(&report, managers)
	l.ScoreRuleSingle(&report, managers)
	l.ScoreRuleUseDev(&report, managers.Main())

	return report, nil
}

func (l *DependenciesLinter) ScoreRuleUse(report *api.Report, managers api.DependencyManagerList) {
	if len(managers) == 0 {
		report.Scores[RuleUse] = 0
		return
	}

	switch {
	case managers.ContainsType(depmanagers.TypePipenv) || managers.ContainsType(depmanagers.TypePoetry):
		report.Scores[RuleUse] = 100
	case managers.ContainsType(depmanagers.TypeRequirementsTxt):
		report.Scores[RuleUse] = 20 // it's better than nothing, but not recommended
		report.Details[RuleUse] = DetailsNoRequirementsTxt
	case managers.ContainsType(depmanagers.TypeSetupPy):
		report.Scores[RuleUse] = 30 // it's better than nothing and slightly better than a requirements.txt, but still not recommended.
		report.Details[RuleUse] = DetailsNoSetupPy
	default:
		report.Scores[RuleUse] = 0
		report.Details[RuleUse] = fmt.Sprintf("Your project is somehow using a dependency manager that mllint recognises, but cannot score: %s.\n\nPlease create an issue on mllint's GitHub :)", types(managers))
	}
}

func (l *DependenciesLinter) ScoreRuleSingle(report *api.Report, managers api.DependencyManagerList) {
	if len(managers) == 0 {
		report.Scores[RuleSingle] = 0
		return
	}

	if len(managers) == 1 {
		report.Scores[RuleSingle] = 100
		return
	}

	report.Scores[RuleSingle] = 0

	details := strings.Builder{}
	details.WriteString(fmt.Sprintf("Your project was found to be using multiple dependency managers: %+v\n\n", types(managers)))
	switch {
	case managers.ContainsAllTypes(depmanagers.TypeRequirementsTxt, depmanagers.TypeSetupPy):
		details.WriteString(DetailsRequirementsTxtSetupPy)
	case managers.ContainsAllTypes(depmanagers.TypeRequirementsTxt, depmanagers.TypePipenv):
		details.WriteString(DetailsRequirementsTxtPipenv)
	case managers.ContainsAllTypes(depmanagers.TypeRequirementsTxt, depmanagers.TypePoetry):
		details.WriteString(DetailsRequirementsTxtPoetry)
	case managers.ContainsAllTypes(depmanagers.TypePipenv, depmanagers.TypeSetupPy):
		details.WriteString(DetailsPipenvSetupPy)
	case managers.ContainsAllTypes(depmanagers.TypePoetry, depmanagers.TypeSetupPy):
		details.WriteString(DetailsPoetrySetupPy)
	default:
		details.WriteString("Pick the one most suited for you, your project and your team, then stick with it.")
	}

	report.Details[RuleSingle] = details.String()
}

func (l *DependenciesLinter) ScoreRuleUseDev(report *api.Report, manager api.DependencyManager) {
	if manager == nil {
		return
	}

	if manager.Type() == depmanagers.TypeRequirementsTxt || manager.Type() == depmanagers.TypeSetupPy {
		report.Scores[RuleUseDev] = 0
		report.Details[RuleUseDev] = "Your project's main dependency manager is a `" + manager.Type().String() + "` file, which doesn't distinguish between regular dependencies and development dependencies."
		return
	}

	deps := manager.Dependencies()
	shouldBeDevDeps, ok := intersect.Hash(deps, ShouldBeDevDependencies).([]interface{})
	if !ok {
		shouldBeDevDeps = []interface{}{}
	}

	missingDevDeps := []interface{}{}
	for _, d := range shouldBeDevDeps {
		if dep := d.(string); !manager.HasDevDependency(dep) {
			missingDevDeps = append(missingDevDeps, dep)
		}
	}

	if len(missingDevDeps) == 0 {
		report.Scores[RuleUseDev] = 100
		return
	}

	// 1 misplaced dependency: 50%, 2 or more misplaced dependencies: 0%
	report.Scores[RuleUseDev] = math.Max(float64(100-50*len(missingDevDeps)), 0)
	report.Details[RuleUseDev] = fmt.Sprint(manager.Type().String(), ` is tracking the following dependencies as regular dependencies, but they should actually be development dependencies.

Please move the following dependencies `, instructionsHowToMovePkgs[manager.Type()], "\n\n", markdowngen.List(missingDevDeps))
}

var instructionsHowToMovePkgs = map[api.DependencyManagerType]string{
	depmanagers.TypePoetry: "from the `dependencies` section to the `dev-dependencies` section in your `pyproject.toml`, then run `poetry lock` to update your lock file.",
	depmanagers.TypePipenv: "from the `packages` section to the `dev-packages` section in your `Pipfile`, then run `pipenv lock` to update your lock file.",
}

func types(managers []api.DependencyManager) []api.DependencyManagerType {
	types := []api.DependencyManagerType{}
	for _, manager := range managers {
		types = append(types, manager.Type())
	}
	return types
}
