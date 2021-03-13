package projectlinters

import (
	"path"

	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/config"
	"gitlab.com/bvobart/mllint/utils"
	"gitlab.com/bvobart/mllint/utils/dvc"
	"gitlab.com/bvobart/mllint/utils/git"
)

const (
	MsgUseDVC = `Your project is not using DVC (Data Version Control). DVC is an open-source version control system for Machine Learning projects.
		DVC is built to help version your data and make ML models shareable and reproducible. It is designed to handle large files, datasets, ML models, and metrics as well as code.
		DVC can also help you manage ML experiments by guaranteeing that all files and metrics will be consistent and in the right place to reproduce the experiments,
		or use it as a baseline for a new iteration.
		> Start using DVC to track your data and ML models. Install DVC and run 'dvc init' to get started.
		> Learn more about DVC and how to use it here: https://dvc.org/doc`

	MsgUseDVCIsInstalled = `Your project is using DVC (Data Version Control), but it doesn't seem to be installed as I cannot find the 'dvc' executable on your PATH.
		> Install and add the 'dvc' package to your dependency manager, e.g. by running 'poetry add dvc'`

	MsgCommitDVCFolder = `Your project is using DVC (Data Version Control), but you haven't committed any of the files in the '.dvc' folder to your Git repository.
		The files in the '.dvc' folder are crucial for DVC to be able to work.
		> Add the files in the '.dvc' folder to your Git repository, then commit and push them:
		>   git add .dvc && git commit -m "Adds .dvc folder" && git push`

	MsgDVCAddRemote = `Your project is using DVC (Data Version Control), but you haven't configured any remotes yet.
		Remotes provide a location to store and share your data and models with the rest of your team.
		> Read more about DVC remotes here: https://dvc.org/doc/command-reference/remote
		> Then add a remote to DVC using 'dvc remote add ...'`

	MsgDVCAddFiles = `Your project is using DVC (Data Version Control), but you haven't added any files to DVC yet.
		It is recommended to version your datasets and ML models 
		> Learn more about adding data to DVC at  https://dvc.org/doc/start/data-versioning
		> and add your datasets and models to DVC:  'dvc add ...'`

	MsgCommitDVCLock = `Your project is using DVC (Data Version Control) and you have a  dvc.lock  file, 
		but you haven't added it to Git yet as is recommended with lock files.
		> Add dvc.lock to Git:  'git add dvc.lock'`
)

const (
	RuleCommitDVCFolder = "commit-dvc-folder"
	RuleCommitDVCLock   = "commit-dvc-lock"
	RuleDVCAddRemote    = "add-remote"
	RuleDVCAddFiles     = "add-files"
	RuleDVCIsInstalled  = "is-installed"
)

// UseDVC is a linting rule that checks whether the project is using DVC.
// Relates to best practice 'Use Versioning for Data, Model, Configurations and Training Script'
// See https://se-ml.github.io/best_practices/02-data_version/
type UseDVC struct{}

func (l *UseDVC) Name() string {
	return "use-dvc"
}

func (l *UseDVC) Rules() []string {
	return []string{"", RuleCommitDVCFolder, RuleCommitDVCLock, RuleDVCAddRemote, RuleDVCAddFiles, RuleDVCIsInstalled}
}

func (l *UseDVC) Configure(_ *config.Config) error {
	return nil
}

func (l *UseDVC) LintProject(projectdir string) ([]api.Issue, error) {
	if !utils.FileExists(path.Join(projectdir, ".dvc", "config")) {
		return []api.Issue{api.NewIssue(l.Name(), "", api.SeverityError, MsgUseDVC)}, nil
	}

	issues := []api.Issue{}

	// Test whether .dvc folder is tracked by Git
	if !git.IsTracking(projectdir, ".dvc") {
		issues = append(issues, api.NewIssue(l.Name(), RuleCommitDVCFolder, api.SeverityError, MsgCommitDVCFolder))
	}

	// Test whether DVC is installed
	if !dvc.IsInstalled() {
		issues = append(issues, api.NewIssue(l.Name(), RuleDVCIsInstalled, api.SeverityError, MsgUseDVCIsInstalled))
	}

	// At this point, if there are issues, return them.
	// There's no point in checking the rest of the linting rules if these basic issues are not resolved.
	if len(issues) > 0 {
		return issues, nil
	}

	// Test whether a remote has been configured: `dvc remote list`
	if len(dvc.Remotes(projectdir)) == 0 {
		issues = append(issues, api.NewIssue(l.Name(), RuleDVCAddRemote, api.SeverityError, MsgDVCAddRemote))
	}

	// Test whether there are any files being tracked with DVC: `dvc list . -R --dvc-only`
	if len(dvc.Files(projectdir)) == 0 {
		issues = append(issues, api.NewIssue(l.Name(), RuleDVCAddFiles, api.SeverityWarning, MsgDVCAddFiles))
	}

	// Check whether the user has committed their dvc.lock file.
	if utils.FileExists(path.Join(projectdir, "dvc.lock")) && !git.IsTracking(projectdir, "dvc.lock") {
		issues = append(issues, api.NewIssue(l.Name(), RuleCommitDVCLock, api.SeverityWarning, MsgCommitDVCLock))
	}

	return issues, nil
}

// Additional DVC linting rules:
// - Check whether the user has defined any DVC pipelines (repo has a dvc.yaml file with non-empty stages) (warning?)
// - something related to experiment tracking?

// File and folder structure
// 	 - Check default file location for data? E.g. always use ./data, so all .dvc files must be in (subdir of) ./data
//	   - From cookiecutter Data Science template (data folder, models folder, src folder?)

// Also interesting
//   - look into Weights & Biases artefacts
//   - MLFlow: experiment tracking, deploy model as API, 'end-to-end', incl deployment. Not very great to use though.
//   - comet - experiment tracking
//   - neptune - experiment tracking
