package projectlinters

import "gitlab.com/bvobart/mllint/api"

type GitLinter struct{}

func (l GitLinter) Name() string {
	return "Git usage checker"
}

func (l GitLinter) LintProject(projectdir string) ([]api.Issue, error) {
	return []api.Issue{}, nil
}
