package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

const ghactionsFolder = ".github/workflows"

type GHActions struct{}

func (_ GHActions) ConfigFile() string {
	return ghactionsFolder
}

func (_ GHActions) Detect(projectdir string) bool {
	workflowsdir := path.Join(projectdir, ghactionsFolder)
	if !utils.FolderExists(workflowsdir) {
		return false
	}

	isEmpty, err := utils.FolderIsEmpty(workflowsdir)
	if err != nil {
		return false
	}

	return !isEmpty
}

func (_ GHActions) Type() ProviderType {
	return TypeGHActions
}
