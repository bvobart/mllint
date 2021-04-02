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
	Details: fmt.Sprintf(`The code of any software project should be tracked in version control software. Git is free, open-source, and the most popular tool for version controlling code, designed to handle anything from small projects to extremely large projects such as the Linux kernel.

Version control software allows you to track changes to your code, easily return to an earlier version and helps to collaborate with other people in developing your project.

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
	Large files should instead be version controlled as Data, e.g. using Git LFS or DVC. See the %s rules in the Version Control category of %s
	
	It is not enough to just remove these large files from your local filesystem, as the files will still exist inside your Git history.
	Instead, see [this StackOverflow answer](https://stackoverflow.com/a/46615578/8059181) to learn how to also remove these files from your project's Git history.`, "`data/`", "`mllint`"),
	Weight: 1,
}
