package testutils

import (
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/commands/mllint"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils"
)

type LinterTestOptions struct {
	conf              *config.Config
	runner            mllint.Runner
	detectPythonFiles bool
	detectDepManagers bool
	detectCQLinters   bool
	usePythonFiles    utils.Filenames
	useDepManagers    api.DependencyManagerList
	useCQLinters      []api.CQLinter
}

func NewOptions() *LinterTestOptions {
	return &LinterTestOptions{}
}

func (opts *LinterTestOptions) DetectPythonFiles() *LinterTestOptions {
	opts.detectPythonFiles = true
	return opts
}

func (opts *LinterTestOptions) DetectDepManagers() *LinterTestOptions {
	opts.detectDepManagers = true
	return opts
}

func (opts *LinterTestOptions) DetectCQLinters() *LinterTestOptions {
	opts.detectCQLinters = true
	return opts
}

func (opts *LinterTestOptions) UsePythonFiles(files utils.Filenames) *LinterTestOptions {
	opts.usePythonFiles = files
	return opts
}

func (opts *LinterTestOptions) UseDepManagers(managers api.DependencyManagerList) *LinterTestOptions {
	opts.useDepManagers = managers
	return opts
}

func (opts *LinterTestOptions) UseCQLinters(linters []api.CQLinter) *LinterTestOptions {
	opts.useCQLinters = linters
	return opts
}

func (opts *LinterTestOptions) WithConfig(conf *config.Config) *LinterTestOptions {
	opts.conf = conf
	return opts
}

func (opts *LinterTestOptions) WithRunner(runner mllint.Runner) *LinterTestOptions {
	opts.runner = runner
	return opts
}
