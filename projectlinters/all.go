package projectlinters

import (
	"gitlab.com/bvobart/mllint/api"
)

func GetAllLinters() api.LinterList {
	return api.LinterList{
		&UseGit{},
		&GitNoBigFiles{},
		&UseDependencyManager{},
	}
}
