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
	// ConfigFile returns the location of the CI provider's configuration file in a project.
	ConfigFile(projectdir string) string

	// Detects whether the project at the given location uses this provider. Checking for config file existance should be enough
	Detect(projectdir string) bool

	// Type of CI provider, i.e. one of ciproviders.Type*
	Type() ProviderType
}

// Detect detects CI providers in the root of the Git repository that the given folder is in.
// Often the Git root dir and the given folder will be the same, but not in the case of monorepo style repos.
// When the dir is not in a Git repo, then it will simply just check the dir.
func Detect(projectdir string) []Provider {
	providers := []Provider{}
	for _, p := range all {
		if p.Detect(projectdir) {
			providers = append(providers, p)
		}
	}
	return providers
}
