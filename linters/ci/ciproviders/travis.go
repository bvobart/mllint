package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

const travisFile = ".travis.yml"

type Travis struct{}

func (_ Travis) ConfigFile() string {
	return travisFile
}

func (_ Travis) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, travisFile))
}

func (_ Travis) Type() ProviderType {
	return TypeTravis
}
