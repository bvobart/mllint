package projectlinters

import (
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/config"
)

// TODO: put this setting in a config
const largeFileThreshold = 10_000_000 // 10 MB

func GetAllLinters() []api.Linter {
	return []api.Linter{
		UseGit{},
		GitNoBigFiles{Threshold: largeFileThreshold},
		UseDependencyManager{},
	}
}

func FilterEnabled(all []api.Linter, conf config.RuleConfig) []api.Linter {
	linters := map[string]api.Linter{}
	for _, linter := range all {
		linters[linter.Name()] = linter
	}

	for _, rule := range conf.Disabled {
		if _, isEnabled := linters[rule]; isEnabled {
			linters[rule] = nil
		}
	}

	enabledLinters := []api.Linter{}
	for _, linter := range linters {
		if linter != nil {
			enabledLinters = append(enabledLinters, linter)
		}
	}

	return enabledLinters
}
