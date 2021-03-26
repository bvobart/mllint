package projectlinters

import (
	"fmt"

	"github.com/dustin/go-humanize"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils/git"
)

const (
	MsgUseGit = `Your project is not using Git. Version control software such as Git allows you to track changes to your code,
		easily return to an earlier version, and help to collaborate with other people in developing your project.
		> Start using Git by running 'git init'`
	MsgNoBigFiles = `File '%s' is being tracked by Git, but it is %s, which is too large (> %s) to comfortably use in a Git repository. 
		> Consider removing this file from the Git index 
		> and versioning it using Data Version Control (DVC), Git LFS, or another tool for versioning large files with Git.`
)

// UseGit is a linting rule that checks whether the project is using Git.
// Relates to best practice 'Use Versioning for Data, Model, Configurations and Training Script'
// See https://se-ml.github.io/best_practices/02-data_version/
type UseGit struct{}

func (l *UseGit) Name() string {
	return "use-git"
}

func (l *UseGit) Rules() []string {
	return []string{""}
}

func (l *UseGit) Configure(_ *config.Config) error {
	return nil
}

func (l *UseGit) LintProject(projectdir string) ([]api.Issue, error) {
	if !git.Detect(projectdir) {
		return []api.Issue{api.NewIssue(l.Name(), "", api.SeverityError, MsgUseGit)}, nil
	}
	return nil, nil
}

// GitNoBigFiles is a linting rule that will check whether there are no big files in the Git repository.
// Relates to best practices of Git usage.
// See https://docs.github.com/en/github/managing-large-files/what-is-my-disk-quota
type GitNoBigFiles struct {
	Threshold uint64
}

func (l *GitNoBigFiles) Name() string {
	return "git-no-big-files"
}

func (l *GitNoBigFiles) Rules() []string {
	return []string{""}
}

func (l *GitNoBigFiles) Configure(conf *config.Config) error {
	l.Threshold = conf.Git.MaxFileSize
	return nil
}

func (l *GitNoBigFiles) LintProject(projectdir string) ([]api.Issue, error) {
	// if this project does not use Git, this linting rule will just crash, so we skip it
	if !git.Detect(projectdir) {
		return nil, nil
	}

	largeFiles, err := git.FindLargeFiles(projectdir, l.Threshold)
	if err != nil {
		return nil, err
	}

	if len(largeFiles) == 0 {
		return nil, nil
	}

	issues := []api.Issue{}
	for _, file := range largeFiles {
		msg := fmt.Sprintf(MsgNoBigFiles, file.Path, humanize.Bytes(file.Size), humanize.Bytes(l.Threshold))
		issues = append(issues, api.NewIssue(l.Name(), "", api.SeverityWarning, msg))
	}

	return issues, nil
}
