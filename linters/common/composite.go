package common

import (
	"fmt"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/config"
	"github.com/hashicorp/go-multierror"
)

func NewCompositeLinter(name string, linters ...api.Linter) mllint.ConfigurableLinterWithRunner {
	rules := []*api.Rule{}
	for _, linter := range linters {
		rules = append(rules, linter.Rules()...)
	}
	return &CompositeLinter{name, linters, rules, nil}
}

type CompositeLinter struct {
	name    string
	linters []api.Linter
	rules   []*api.Rule
	runner  mllint.Runner
}

func (l *CompositeLinter) Name() string {
	return l.name
}

func (l *CompositeLinter) Rules() []*api.Rule {
	return l.rules
}

func (l *CompositeLinter) SetRunner(r mllint.Runner) {
	l.runner = r
}

func (l *CompositeLinter) Configure(conf *config.Config) error {
	for _, linter := range l.linters {
		configurable, ok := linter.(api.Configurable)
		if ok {
			if err := configurable.Configure(conf); err != nil {
				return fmt.Errorf("configuration error in linter '%s': %w", linter.Name(), err)
			}
		}
	}
	return nil
}

func (l *CompositeLinter) LintProject(project api.Project) (api.Report, error) {
	tasks := make([]*mllint.RunnerTask, len(l.linters))
	for i, linter := range l.linters {
		displayName := l.name + " - " + linter.Name()
		tasks[i] = l.runner.RunLinter(fmt.Sprint(i), linter, project, mllint.DisplayName(displayName))
	}

	var err *multierror.Error
	reports := make([]api.Report, 0, len(tasks))
	mllint.ForEachTask(l.runner.CollectTasks(tasks...), func(task *mllint.RunnerTask, result mllint.LinterResult) {
		if result.Err != nil {
			err = multierror.Append(err, fmt.Errorf("linting error in linter '%s': %w", task.Linter.Name(), result.Err))
		}

		reports = append(reports, result.Report)
	})

	return api.MergeReports(api.NewReport(), reports...), err.ErrorOrNil()
}
