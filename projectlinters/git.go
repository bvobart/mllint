package projectlinters

import (
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/utils/git"
)

type GitLinter struct{}

func (l GitLinter) Name() string {
	return "Git usage checker"
}

func (l GitLinter) LintProject(projectdir string) ([]api.Issue, error) {
	if !git.Detect(projectdir) {
		issue := api.NewIssue(api.SeverityError, "Your project is not using Git.") // TODO: expand this message somewhat more
		return []api.Issue{issue}, nil
	}
	return []api.Issue{}, nil
}
