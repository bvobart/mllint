package projectlinters

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/utils/git"
)

const (
	MsgUseGit     = "Your project is not using Git. Version control software such as Git allows you to track changes to your code, easily return to an earlier version, and help to collaborate with other people in developing your project."
	MsgNoBigFiles = "File '%s' is being tracked by Git, but it is %s, which is too large (> %s) to comfortably use in a Git repository. Consider removing this file from the Git index and versioning it using Git LFS, Data Version Control (DVC), or another tool for versioning large files with Git."
)

// UseGit is a linting rule that checks whether the project is using Git.
type UseGit struct{}

func (l UseGit) Name() string {
	return "use-git"
}

func (l UseGit) LintProject(projectdir string) ([]api.Issue, error) {
	if !git.Detect(projectdir) {
		return []api.Issue{api.NewIssue(l.Name(), api.SeverityError, MsgUseGit)}, nil
	}
	return nil, nil
}

// GitNoBigFiles is a linting rule that will check whether there are no big files in the Git repository.
type GitNoBigFiles struct {
	Threshold uint64
}

func (l GitNoBigFiles) Name() string {
	return "git-no-big-files"
}

func (l GitNoBigFiles) LintProject(projectdir string) ([]api.Issue, error) {
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
		msg := fmt.Sprintf(MsgNoBigFiles, file.Path, humanize.Bytes(file.Size), humanize.Bytes(largeFileThreshold))
		issues = append(issues, api.NewIssue(l.Name(), api.SeverityWarning, msg))
	}

	return issues, nil
}
