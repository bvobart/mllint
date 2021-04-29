package codequality

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func NewLinter() api.ConfigurableLinter {
	return &CQLinter{}
}

type CQLinter struct {
	Linters []api.CQLinter
}

func (l *CQLinter) Name() string {
	return categories.CodeQuality.Name
}

func (l *CQLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleUseLinters}
}

func (l *CQLinter) Configure(conf *config.Config) (err error) {
	l.Linters, err = cqlinters.FromConfig(conf.CodeQuality)
	return err
}

func (l *CQLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	detectedLinters := project.CQLinters
	desiredLinters := l.Linters
	if len(desiredLinters) == 0 {
		return report, nil
	}

	missingLinters := findMissing(desiredLinters, detectedLinters)
	if len(missingLinters) == 0 {
		report.Scores[RuleUseLinters] = 100
		report.Details[RuleUseLinters] = "All linters detected!\n\n" + markdowngen.List(asInterfaceList(desiredLinters))
	} else {
		report.Scores[RuleUseLinters] = 100 - 100*float64(len(missingLinters)/len(desiredLinters))
		report.Details[RuleUseLinters] = "Your project should employ the following linters to help you measure the quality of your code:\n\n" + markdowngen.List(asInterfaceList(missingLinters))
	}

	for _, desiredLinter := range desiredLinters {
		// TODO: desiredLinter.Run()
		// then gather linting issues
		// find a way to score the specific rule for this CQlinter
		_ = desiredLinter
	}

	// TODO: for each configured / desired linter:
	// 					check whether there is a configuration for it in the repository?
	//          run linter
	//          gather issues
	//          count to gain score
	// may want separate linter objects for this.

	return report, nil
}

func contains(linters []api.CQLinter, target api.CQLinter) bool {
	for _, l := range linters {
		if l == target {
			return true
		}
	}
	return false
}

func findMissing(desired []api.CQLinter, detected []api.CQLinter) []api.CQLinter {
	missing := []api.CQLinter{}
	for _, desiredLinter := range desired {
		if !contains(detected, desiredLinter) {
			missing = append(missing, desiredLinter)
		}
	}
	return missing
}

func asInterfaceList(list []api.CQLinter) []interface{} {
	res := make([]interface{}, len(list))
	for i := range list {
		res[i] = list[i]
	}
	return res
}
