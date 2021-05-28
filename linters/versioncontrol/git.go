package versioncontrol

import (
	"fmt"
	"math"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/setools/git"
)

const penaltyPerLargeFile = 25 // percent

type GitLinter struct {
	MaxFileSize uint64
}

func (l *GitLinter) Name() string {
	return "Git"
}

func (l *GitLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleGit, &RuleGitNoBigFiles}
}

func (l *GitLinter) Configure(conf *config.Config) error {
	l.MaxFileSize = conf.Git.MaxFileSize
	return nil
}

func (l *GitLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	report.Scores[RuleGit] = 100
	if !git.Detect(project.Dir) {
		report.Scores[RuleGit] = 0
		return report, nil
	}

	largeFiles, err := git.FindLargeFilesInHistory(project.Dir, l.MaxFileSize)
	if err != nil {
		return api.Report{}, err
	}

	report.Scores[RuleGitNoBigFiles] = math.Max(float64(100-penaltyPerLargeFile*len(largeFiles)), 0)
	if len(largeFiles) > 0 {
		report.Details[RuleGitNoBigFiles] = l.buildDetails(largeFiles)
	}

	return report, nil
}

func (l *GitLinter) buildDetails(largeFiles []git.FileSize) string {
	details := strings.Builder{}
	details.WriteString(fmt.Sprintf("Your project's Git history contains the following files that are larger than %s:\n", humanize.Bytes(l.MaxFileSize)))
	for _, file := range largeFiles {
		details.WriteString(fmt.Sprintf("- **%s** - commit `%s` - `%s`\n", humanize.Bytes(file.Size), file.CommitHash, file.Path))
	}
	details.WriteString(`
These files may not necessarily be in your project right now, but they are still stored inside your project's Git history.
Thus, whenever your project is downloaded (with ` + "`git clone`" + `), all these unnecessary files have to be downloaded as well.

See [this StackOverflow answer](https://stackoverflow.com/a/46615578/8059181) to learn how to remove these files from your project's Git history.`)
	return details.String()
}
