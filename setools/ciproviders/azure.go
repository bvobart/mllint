package ciproviders

import (
	"path"

	"github.com/bvobart/mllint/utils"
)

const azureFile = "azure-pipelines.yml"

type Azure struct{}

func (_ Azure) ConfigFile() string {
	return azureFile
}

func (_ Azure) Detect(projectdir string) bool {
	return utils.FileExists(path.Join(projectdir, azureFile))
}

func (_ Azure) Type() ProviderType {
	return TypeAzure
}
