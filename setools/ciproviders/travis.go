package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/setools/git"
	"github.com/bvobart/mllint/utils"
)

const travisFile = ".travis.yml"

type Travis struct{}

func (_ Travis) ConfigFile(projectdir string) string {
	return path.Join(git.GetGitRoot(projectdir), travisFile)
}

func (_ Travis) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(git.GetGitRoot(projectdir), travisFile))
}

func (_ Travis) Type() ProviderType {
	return TypeTravis
}
