package versioncontrol

import (
	"path"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/dvc"
	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils"
)

// DVCLinter is a linter that checks whether the project is using DVC and using it correctly and effectively.
// Relates to best practice 'Use Versioning for Data, Model, Configurations and Training Script'
// See https://se-ml.github.io/best_practices/02-data_version/
type DVCLinter struct{}

func (l *DVCLinter) Name() string {
	return "Data"
}

func (l *DVCLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleDVC, &RuleDVCIsInstalled, &RuleCommitDVCFolder, &RuleDVCHasRemote, &RuleDVCHasFiles, &RuleCommitDVCLock}
}

func (l *DVCLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()
	report.Scores[RuleDVC] = 0
	report.Scores[RuleDVCIsInstalled] = 0
	report.Scores[RuleCommitDVCFolder] = 0
	report.Scores[RuleDVCHasRemote] = 0
	report.Scores[RuleDVCHasFiles] = 0

	// Test whether 'dvc init' was run by checking whether a .dvc/config exists.
	if utils.FileExists(path.Join(project.Dir, ".dvc", "config")) {
		report.Scores[RuleDVC] = 100
	} else {
		return report, nil
	}

	if git.IsTracking(project.Dir, ".dvc") {
		report.Scores[RuleCommitDVCFolder] = 100
	}

	// Test whether DVC is installed. If it is not, then the other rules below cannot be checked, so we return.
	if dvc.IsInstalled() {
		report.Scores[RuleDVCIsInstalled] = 100
	} else {
		return report, nil
	}

	// Test whether a remote has been configured: `dvc remote list`
	if !RuleDVCHasRemote.Disabled && len(dvc.Remotes(project.Dir)) > 0 {
		report.Scores[RuleDVCHasRemote] = 100
	}

	// Test whether there are any files being tracked with DVC: `dvc list . -R --dvc-only`
	if !RuleDVCHasFiles.Disabled && len(dvc.Files(project.Dir)) > 0 {
		report.Scores[RuleDVCHasFiles] = 100
	}

	// Check whether the user has committed their dvc.lock file.
	if utils.FileExists(path.Join(project.Dir, "dvc.lock")) {
		report.Scores[RuleCommitDVCLock] = 0
		if git.IsTracking(project.Dir, "dvc.lock") {
			report.Scores[RuleCommitDVCLock] = 100
		}
	}

	return report, nil
}

// Ideas for future DVC linting rules:
// - Check whether the user has defined any DVC pipelines (repo has a dvc.yaml file with non-empty stages)
// - Check what kinds of pipelines and stages the user has defined, check if there's a cleaning stage, training stage, testing stage, etc.
// - something related to experiment tracking?

// File and folder structure
// 	 - Check default file location for data? E.g. always use ./data, so all .dvc files must be in (subdir of) ./data
//	   - From cookiecutter Data Science template (data folder, models folder, src folder?)

// Also interesting
//   - look into Weights & Biases artefacts
//   - MLFlow: experiment tracking, deploy model as API, 'end-to-end', incl deployment. Not very great to use though.
//   - comet - experiment tracking
//   - neptune - experiment tracking
