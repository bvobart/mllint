package ciproviders_test

import (
	"os"
	"path"
	"testing"

	"github.com/bvobart/mllint/setools/ciproviders"
	"github.com/bvobart/mllint/setools/git"
	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	dir := "."
	providers := ciproviders.Detect(dir)
	require.Equal(t, []ciproviders.Provider{ciproviders.GHActions{}}, providers)

	dir = os.TempDir()
	providers = ciproviders.Detect(dir)
	require.Equal(t, []ciproviders.Provider{}, providers)
}

func TestConfigFile(t *testing.T) {
	dir := "."
	provider := ciproviders.GHActions{}
	configDir := provider.ConfigFile(dir)
	require.Equal(t, path.Join(git.GetGitRoot(dir), ".github", "workflows"), configDir)
}
