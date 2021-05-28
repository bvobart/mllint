package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils"
)

const gitlabFile = ".gitlab-ci.yml"

type Gitlab struct{}

func (_ Gitlab) ConfigFile(projectdir string) string {
	return path.Join(git.GetGitRoot(projectdir), gitlabFile)
}

func (_ Gitlab) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(git.GetGitRoot(projectdir), gitlabFile))
}

func (_ Gitlab) Type() ProviderType {
	return TypeGitlab
}
