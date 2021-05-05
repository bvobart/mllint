package ciproviders

const (
	TypeAzure     ProviderType = "Azure DevOps"
	TypeGHActions ProviderType = "GitHub Actions"
	TypeGitlab    ProviderType = "GitLab CI"
	TypeTravis    ProviderType = "Travis CI"
)

var (
	azure     Provider = Azure{}
	ghActions Provider = GHActions{}
	gitlab    Provider = Gitlab{}
	travis    Provider = Travis{}
)

var all = []Provider{azure, ghActions, gitlab, travis}

type ProviderType string
type Provider interface {
	// ConfigFile returns the location of the CI provider's configuration file, relative to the project's root.
	ConfigFile() string

	// Detects whether the project at the given location uses this provider. Checking for config file existance should be enough
	Detect(projectdir string) bool

	// Type of CI provider, i.e. one of ciproviders.Type*
	Type() ProviderType
}

func Detect(projectdir string) []Provider {
	providers := []Provider{}
	for _, p := range all {
		if p.Detect(projectdir) {
			providers = append(providers, p)
		}
	}
	return providers
}
