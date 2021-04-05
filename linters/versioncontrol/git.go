package versioncontrol

import (
	"fmt"
	"math"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/versioncontrol/git"
)

const penaltyPerLargeFile = 25 // percent

type GitLinter struct {
	MaxFileSize uint64
}

// Name returns "Code". That might be strange since the struct is called the GitLinter,
// but it is really linter whether / how the project is version controlling code,
// for which obviously by far most popular solution nowadays is Git.
func (l *GitLinter) Name() string {
	return "Code"
}

func (l *GitLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleGit, &RuleGitNoBigFiles}
}

func (l *GitLinter) Configure(conf *config.Config) error {
	l.MaxFileSize = conf.Git.MaxFileSize
	return nil
}

func (l *GitLinter) LintProject(projectdir string) (api.Report, error) {
	report := api.NewReport()

	report.Scores[RuleGit] = 100
	if !git.Detect(projectdir) {
		report.Scores[RuleGit] = 0
		return report, nil
	}

	largeFiles, err := git.FindLargeFiles(projectdir, l.MaxFileSize)
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
	details.WriteString(fmt.Sprintf("Your project is tracking the following files that are larger than %s:\n", humanize.Bytes(l.MaxFileSize)))
	for _, file := range largeFiles {
		details.WriteString(fmt.Sprintf("- **%s**  %s\n", humanize.Bytes(file.Size), file.Path))
	}
	return details.String()
}
