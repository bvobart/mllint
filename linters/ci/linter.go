package ci

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/ciproviders"
	"github.com/bvobart/mllint/setools/git"
)

func NewLinter() api.Linter {
	return &CILinter{}
}

type CILinter struct{}

func (l *CILinter) Name() string {
	return "Continuous Integration (CI)"
}

func (l *CILinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleUseCI}
}

func (l *CILinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	report.Scores[RuleUseCI] = 0

	providers := ciproviders.Detect(project.Dir)
	if len(providers) > 0 {
		report.Scores[RuleUseCI] = 100
	}

	for _, provider := range providers {
		// if the repo is not tracking the CI config file, then they're not really using CI,
		// they're merely trying to define it, which is at least a step in the right direction.
		if !git.IsTracking(project.Dir, provider.ConfigFile()) {
			report.Scores[RuleUseCI] = 25
		}
	}

	return report, nil
}
