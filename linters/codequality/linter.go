package codequality

import (
	"github.com/hashicorp/go-multierror"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/codequality/bandit"
	"github.com/bvobart/mllint/linters/codequality/black"
	"github.com/bvobart/mllint/linters/codequality/isort"
	"github.com/bvobart/mllint/linters/codequality/mypy"
	"github.com/bvobart/mllint/linters/codequality/pylint"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils/markdowngen"
)

type pair struct {
	Linter api.Linter
	Type   api.CQLinterType
}

var all = []pair{
	{pylint.NewLinter(), cqlinters.TypePylint},
	{mypy.NewLinter(), cqlinters.TypeMypy},
	{black.NewLinter(), cqlinters.TypeBlack},
	{isort.NewLinter(), cqlinters.TypeISort},
	{bandit.NewLinter(), cqlinters.TypeBandit},
}

func toMap(all []pair) map[api.CQLinterType]api.Linter {
	m := make(map[api.CQLinterType]api.Linter, len(all))
	for _, l := range all {
		m[l.Type] = l.Linter
	}
	return m
}

var sublinters = toMap(all)

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
	rules := []*api.Rule{&RuleUseLinters, &RuleLintersInstalled}
	for _, l := range all {
		rules = append(rules, l.Linter.Rules()...)
	}
	return rules
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
		report.Scores[RuleUseLinters] = 100 * (1 - float64(len(missingLinters)/len(desiredLinters)))
		report.Details[RuleUseLinters] = "Your project should employ the following linters to help you measure the quality of your code:\n\n" + markdowngen.List(asInterfaceList(missingLinters))
	}

	notInstalledLinters := findNotInstalled(desiredLinters)
	if len(notInstalledLinters) == 0 {
		report.Scores[RuleLintersInstalled] = 100
	} else {
		report.Scores[RuleLintersInstalled] = 100 * (1 - float64(len(notInstalledLinters)/len(desiredLinters)))
		report.Details[RuleLintersInstalled] = "The following linters were not installed, so we could not analyse what they had to say about your project:\n\n" + markdowngen.List(asInterfaceList(notInstalledLinters))
	}

	// TODO: run all these linters in parallel

	var multiErr *multierror.Error
	subReports := []api.Report{}
	for _, desiredLinter := range desiredLinters {
		if !desiredLinter.IsInstalled() {
			continue
		}

		mlLinter, ok := sublinters[desiredLinter.Type()]
		if !ok {
			continue
		}

		subReport, err := mlLinter.LintProject(project)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}

		subReports = append(subReports, subReport)
	}

	return api.MergeReports(report, subReports...), multiErr.ErrorOrNil()
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

func findNotInstalled(desired []api.CQLinter) []api.CQLinter {
	notInstalled := []api.CQLinter{}
	for _, desiredLinter := range desired {
		if !desiredLinter.IsInstalled() {
			notInstalled = append(notInstalled, desiredLinter)
		}
	}
	return notInstalled
}

func asInterfaceList(list []api.CQLinter) []interface{} {
	res := make([]interface{}, len(list))
	for i := range list {
		res[i] = list[i]
	}
	return res
}
