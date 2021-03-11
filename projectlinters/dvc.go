package projectlinters

import (
	"path"

	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/config"
	"gitlab.com/bvobart/mllint/utils"
)

const (
	MsgUseDVC = `Your project is not using DVC (Data Version Control). DVC is an open-source version control system for Machine Learning projects.
		DVC is built to make ML models shareable and reproducible. It is designed to handle large files, datasets, ML models, and metrics as well as code.
		DVC also helps you manage ML experiments by guaranteeing that all files and metrics will be consistent and in the right place to reproduce the experiment
		or use it as a baseline for a new iteration
		> Start using DVC to track your datasets and ML models. Install DVC and run 'dvc init' to get started.
		> Learn more about DVC here: https://dvc.org/doc`
)

// UseDVC is a linting rule that checks whether the project is using DVC.
// Relates to best practice 'Use Versioning for Data, Model, Configurations and Training Script'
// See https://se-ml.github.io/best_practices/02-data_version/
type UseDVC struct{}

func (l *UseDVC) Name() string {
	return "use-dvc"
}

func (l *UseDVC) Rules() []string {
	return []string{""}
}

func (l *UseDVC) Configure(_ *config.Config) error {
	return nil
}

func (l *UseDVC) LintProject(projectdir string) ([]api.Issue, error) {
	if utils.FileExists(path.Join(projectdir, ".dvc", "config")) {
		return nil, nil
	}
	return []api.Issue{api.NewIssue(l.Name(), "", api.SeverityError, MsgUseDVC)}, nil
}
