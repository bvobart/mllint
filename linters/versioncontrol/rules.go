package versioncontrol

import (
	"fmt"

	"github.com/bvobart/mllint/api"
)

// RuleGit is a linting rule to check whether the project is using Git.
// Relates to best practice 'Use Versioning for Data, Model, Configurations and Training Script'
// See https://se-ml.github.io/best_practices/02-data_version/
var RuleGit = api.Rule{
	Slug: "git",
	Name: "Project uses Git",
	Details: fmt.Sprintf(`The code of any software project should be tracked in version control software.
Git is the most widely-used, most popular, free and open-source version controlling tool, designed to handle anything from small projects to extremely large projects such as the Linux kernel.

To start using Git, run %s in a terminal at the root of your project. See also [Git's documentation](https://git-scm.com/doc) for tutorials on how to work with Git.`, "`git init`"),
	Weight: 1,
}

// RuleGitNoBigFiles is a linting rule to check whether there are no big files being tracked in the Git repository.
// Relates to best practices of Git usage.
// See https://docs.github.com/en/github/managing-large-files/what-is-my-disk-quota
var RuleGitNoBigFiles = api.Rule{
	Slug: "git-no-big-files",
	Name: "Project should not use Git to track large files",
	Details: fmt.Sprintf(`Git is great for version controlling small, textual files, but not for binary or large files.
Tracking large files directly with Git adds bloat to your repository's Git history, which needs to be downloaded every time your project is cloned.
Large files should instead be version controlled as Data, e.g. using Git LFS or DVC. See the %s rules in the Version Control category of %s
	
It is not enough to just remove these large files from your local filesystem, as the files will still exist inside your Git history.
Instead, see [this StackOverflow answer](https://stackoverflow.com/a/46615578/8059181) to learn how to also remove these files from your project's Git history.
`, "`data/`", "`mllint`"),
	Weight: 1,
}

//------------------------------------------------------------------------------------------

var RuleDVC = api.Rule{
	Slug: "dvc",
	Name: "Project uses Data Version Control (DVC)",
	Details: fmt.Sprintf(`Similar to code, data should also be version controlled. However, version controlling data cannot be done with Git directly,
as Git is not designed to deal with large and / or binary files. Tracking large files directly with Git adds bloat to your repository's Git history, 
which needs to be downloaded every time your project is cloned.

For properly version controlling Data in ML projects, mllint recommends using [Data Version Control (DVC)](https://dvc.org/).
DVC is an open-source version control system for Machine Learning projects. DVC is built to help version your data and make ML models shareable and reproducible. 
It is designed to handle large files, datasets, ML models, and metrics as well as code.
DVC can also help you manage ML experiments by guaranteeing that all files and metrics will be consistent and in the right place to reproduce the experiments,
or use it as a baseline for a new iteration.

Install DVC (e.g. using %s) and run %s in your terminal to get started with DVC.

To learn more about DVC and how to use it, feel free to check out DVC's documentation and tutorials from these links:
- [DVC Documentation](https://dvc.org/doc)
- [DVC Getting Started](https://dvc.org/doc/start)
- [DVC Installation](https://dvc.org/doc/install)
- [DVC Use Cases](https://dvc.org/doc/use-cases)
- [DVC User Guide](https://dvc.org/doc/user-guide)

Or if you prefer learning from watching videos, DVC has a YouTube channel with several short, useful and informative videos.
- YouTube Channel: [DVCorg](https://www.youtube.com/channel/UC37rp97Go-xIX3aNFVHhXfQ)
- YouTube Video: [Version Control for Data Science Explained in 5 Minutes](https://www.youtube.com/watch?v=UbL7VUpv1Bs)
- YouTube Playlist: [DVC Basics](https://www.youtube.com/playlist?list=PL7WG7YrwYcnDb0qdPl9-KEStsL-3oaEjg)
`, "`poetry add dvc`", "`dvc init`"),
	Weight: 1,
}

var RuleDVCIsInstalled = api.Rule{
	Slug: "dvc-is-installed",
	Name: "DVC is installed",
	Details: fmt.Sprintf(`To be able to use DVC, it must be installed correctly. If you're seeing this as part of an mllint report,
then it means that mllint was unable to find 'dvc' on your PATH. This could either indicate that DVC is not installed in your project,
or it is not included on your path.

See DVC's [installation instructions](https://dvc.org/doc/install) to learn more about installing DVC,
or simply add it to your project as a Pip package, e.g. using %s`, "`poetry add dvc`"),
	Weight: 1,
}

var RuleCommitDVCFolder = api.Rule{
	Slug: "commit-dvc-folder",
	Name: "Folder '.dvc' should be committed to Git",
	Details: fmt.Sprintf(`DVC uses the '.dvc' folder to keep records of and information about all your DVC-tracked files and where they are hosted.
This folder *must* be committed to your Git repository in order to work with DVC correctly.
Learn more about the .dvc directory [here](https://dvc.org/doc/user-guide/project-structure/internal-files).

If you're seeing this in a report, then your project's Git repository is not tracking the '.dvc' folder.
To fix this, you may use the following commands:
%s
`, "```console\ngit add .dvc\ngit commit -m 'Adds .dvc folder for Data Version Control'\ngit push\n```"),
	Weight: 1,
}

var RuleCommitDVCLock = api.Rule{
	Slug: "commit-dvc-lock",
	Name: "File 'dvc.lock' should be committed to Git",
	Details: fmt.Sprintf(`While using DVC to define pipelines in a 'dvc.yaml' file, DVC maintains a 'dvc.lock' file
to record the state of your pipeline(s) and help track its outputs. As with any .lock file, it is highly
recommended to commit your 'dvc.lock' to your project's Git repository. Learn more about dvc.lock files
[here](https://dvc.org/doc/user-guide/project-structure/pipelines-files#dvclock-file).

If you're seeing this in a report, then your project contains a 'dvc.lock' file, but it has not been added to Git.
To add and commit dvc.lock to Git, you may use the following commands:
%s
`, "```console\ngit add dvc.lock\ngit commit -m 'Adds dvc.lock file'\ngit push\n```"),
	Weight: 1,
}

var RuleDVCHasRemote = api.Rule{
	Slug: "dvc-has-remote",
	Name: "DVC should have at least one remote data storage configured",
	Details: fmt.Sprintf(`To share your DVC-tracked data with your colleagues and also allow them to interact with your data,
DVC should have at least one remote storage configured. If you're seeing this in a report, your project currently has none.

Learn more about DVC remotes [here](https://dvc.org/doc/command-reference/remote), then pick your desired remote storage solution,
check the documetation for [adding remotes](https://dvc.org/doc/command-reference/remote/add) and add it as your default remote 
to DVC using %s.
`, "`dvc remote add -d <name> <url>`"),
	Weight: 1,
}

var RuleDVCHasFiles = api.Rule{
	Slug: "dvc-has-files",
	Name: "DVC should be tracking at least one data file",
	Details: fmt.Sprintf(`Using DVC entails tracking changes to your data and models with DVC. If you're seeing this in a report,
your project is using DVC, but it is currently not tracking any files with it.
Learn more about [getting started with data versioning](https://dvc.org/doc/start/data-and-model-versioning) with DVC,
or the [%s](https://dvc.org/doc/command-reference/add) command.

Then, add your datasets and models to DVC by running the command %s

_Tip: Under the hood, mllint uses the command %s in order to see which files DVC is tracking._
`, "`dvc add`", "`dvc add <files>`", "`dvc list . -R --dvc-only`"),
	Weight: 1,
}
