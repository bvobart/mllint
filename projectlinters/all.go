package projectlinters

import "gitlab.com/bvobart/mllint/api"

func GetAllLinters() []api.Linter {
	return []api.Linter{
		GitLinter{},
	}
}
