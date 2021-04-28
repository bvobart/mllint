package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

const gitlabFile = ".gitlab-ci.yml"

type Gitlab struct{}

func (_ Gitlab) ConfigFile() string {
	return gitlabFile
}

func (_ Gitlab) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, gitlabFile))
}

func (_ Gitlab) Type() ProviderType {
	return TypeGitlab
}
