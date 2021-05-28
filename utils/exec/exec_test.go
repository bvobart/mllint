package exec_test

import (
	"testing"

	"github.com/bvobart/mllint/utils/exec"
	"github.com/stretchr/testify/require"
)

func TestDefaultLookPath(t *testing.T) {
	_, err := exec.DefaultLookPath("ls")
	require.NoError(t, err)

	_, err = exec.DefaultLookPath("grep")
	require.NoError(t, err)

	_, err = exec.DefaultLookPath("wc")
	require.NoError(t, err)

	_, err = exec.DefaultLookPath("go")
	require.NoError(t, err)

	_, err = exec.DefaultLookPath("somethingweird")
	require.Error(t, err)
}

func TestDefaultCommandOutput(t *testing.T) {
	output, err := exec.CommandOutput(".", "ls", "-a")
	require.NoError(t, err)
	require.Equal(t, []byte(".\n..\nexec.go\nexec_test.go\nmockexec\n"), output)
}

func TestDefaultCommandCombinedOutput(t *testing.T) {
	output, err := exec.CommandCombinedOutput(".", "ls", "-a")
	require.NoError(t, err)
	require.Equal(t, []byte(".\n..\nexec.go\nexec_test.go\nmockexec\n"), output)
}

func TestDefaultPipelineOutput(t *testing.T) {
	output, err := exec.DefaultPipelineOutput(".", [][]string{
		{"ls", "-al"},
		{"grep", "exec"},
		{"wc", "-l"},
	}...)
	require.NoError(t, err)
	require.Equal(t, []byte("3\n"), output)
}
