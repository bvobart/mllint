package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils"
)

const ghactionsFolder = ".github/workflows"

type GHActions struct{}

func (_ GHActions) ConfigFile(projectdir string) string {
	return path.Join(git.GetGitRoot(projectdir), ghactionsFolder)
}

func (_ GHActions) Detect(projectdir string) bool {
	workflowsdir := path.Join(git.GetGitRoot(projectdir), ghactionsFolder)
	if !utils.FolderExists(workflowsdir) {
		return false
	}

	isEmpty, err := utils.FolderIsEmpty(workflowsdir)
	return err == nil && !isEmpty
}

func (_ GHActions) Type() ProviderType {
	return TypeGHActions
}
