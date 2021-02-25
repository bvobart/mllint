package projectlinters

import "gitlab.com/bvobart/mllint/api"

// TODO: put this setting in a config
const largeFileThreshold = 10_000_000 // 10 MB

func GetAllLinters() []api.Linter {
	return []api.Linter{
		UseGit{},
		GitNoBigFiles{Threshold: largeFileThreshold},
	}
}
