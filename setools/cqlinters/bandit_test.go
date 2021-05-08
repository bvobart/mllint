package cqlinters_test

import (
	"errors"
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
	"github.com/bvobart/mllint/utils/exec/mockexec"
	"github.com/stretchr/testify/require"
)

func TestBandit(t *testing.T) {
	l := cqlinters.Bandit{}
	require.Equal(t, cqlinters.TypeBandit, l.Type())
	require.Equal(t, "Bandit", l.String())
	require.Equal(t, "bandit", l.DependencyName())

	exec.LookPath = mockexec.ExpectLookPath(t, "bandit").ToBeError()
	require.False(t, l.IsInstalled())
	exec.LookPath = mockexec.ExpectLookPath(t, "bandit").ToBeFound()
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath

	project := api.Project{Dir: "."}
	require.False(t, l.IsConfigured(project))
	project.Dir = "test-resources"
	require.True(t, l.IsConfigured(project))
}

func TestBanditRun(t *testing.T) {
	l := cqlinters.Bandit{}
	t.Run("EmptyProject", func(t *testing.T) {
		results, err := l.Run(api.Project{})
		require.NoError(t, err)
		require.Equal(t, results, []api.CQLinterResult{})
	})

	t.Run("NormalProject+String", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = mockexec.ExpectCommand(t).Dir(project.Dir).
			CommandName("bandit").CommandArgs("-f", "yaml", "-x", "test/.venv,test/venv", "-r", project.Dir).
			ToOutput([]byte(testBanditOutput), errors.New("bandit always exits with an error when there are messages"))

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 4)
		for i, result := range results {
			require.IsType(t, cqlinters.BanditMessage{}, result)
			require.Equal(t, expectedBanditOutput[i], result)
		}
	})
}

const testBanditOutput = `errors: []
generated_at: '2021-05-08T15:18:24Z'
metrics:
  ./build/build/lib/mllint/__init__.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./build/build/lib/mllint/__main__.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 2
    nosec: 0
  ./build/build/lib/mllint/cli.py:
    CONFIDENCE.HIGH: 2.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 2.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 10
    nosec: 0
  ./build/mllint/__init__.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./build/mllint/__main__.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 2
    nosec: 0
  ./build/mllint/cli.py:
    CONFIDENCE.HIGH: 2.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 2.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 12
    nosec: 0
  ./build/setup.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 126
    nosec: 0
  ./linters/dependencymgmt/test-resources/multiple/pipenv+setuppy/setup.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./linters/dependencymgmt/test-resources/multiple/poetry+setuppy/setup.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./linters/dependencymgmt/test-resources/multiple/requirementstxt+setuppy/setup.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./linters/dependencymgmt/test-resources/setuppy/setup.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 0
    nosec: 0
  ./utils/test-resources/python-files/some_other_script.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 1
    nosec: 0
  ./utils/test-resources/python-files/some_script.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 12
    nosec: 0
  ./utils/test-resources/python-files/subfolder/yet_another_script.py:
    CONFIDENCE.HIGH: 0.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 0.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 1
    nosec: 0
  _totals:
    CONFIDENCE.HIGH: 4.0
    CONFIDENCE.LOW: 0.0
    CONFIDENCE.MEDIUM: 0.0
    CONFIDENCE.UNDEFINED: 0.0
    SEVERITY.HIGH: 0.0
    SEVERITY.LOW: 4.0
    SEVERITY.MEDIUM: 0.0
    SEVERITY.UNDEFINED: 0.0
    loc: 166
    nosec: 0
results:
- code: 1 import subprocess\n2 import sys\n3 import os\n
  filename: ./build/build/lib/mllint/cli.py
  issue_confidence: HIGH
  issue_severity: LOW
  issue_text: Consider possible security implications associated with subprocess module.
  line_number: 1
  line_range:
  - 1
  more_info: https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b404-import-subprocess
  test_id: B404
  test_name: blacklist
- code: '8   os.chmod(mllint_exe, os.stat(mllint_exe).st_mode | 0o111) # Ensures mllint-exe
    is executable, equivalent to chmod +x\n9   return subprocess.run([mllint_exe]
    + sys.argv[1:]).returncode\n10 \n'
  filename: ./build/build/lib/mllint/cli.py
  issue_confidence: HIGH
  issue_severity: LOW
  issue_text: subprocess call - check for execution of untrusted input.
  line_number: 9
  line_range:
  - 9
  more_info: https://bandit.readthedocs.io/en/latest/plugins/b603_subprocess_without_shell_equals_true.html
  test_id: B603
  test_name: subprocess_without_shell_equals_true
- code: 1 import subprocess\n2 import sys\n3 import os\n
  filename: ./build/mllint/cli.py
  issue_confidence: HIGH
  issue_severity: LOW
  issue_text: Consider possible security implications associated with subprocess module.
  line_number: 1
  line_range:
  - 1
  more_info: https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b404-import-subprocess
  test_id: B404
  test_name: blacklist
- code: '9   os.chmod(mllint_exe, os.stat(mllint_exe).st_mode | 0o111) # Ensures mllint-exe
    is executable, equivalent to chmod +x\n10   return subprocess.run([mllint_exe]
    + sys.argv[1:], check=False).returncode\n11 \n'
  filename: ./build/mllint/cli.py
  issue_confidence: HIGH
  issue_severity: LOW
  issue_text: subprocess call - check for execution of untrusted input.
  line_number: 10
  line_range:
  - 10
  more_info: https://bandit.readthedocs.io/en/latest/plugins/b603_subprocess_without_shell_equals_true.html
  test_id: B603
  test_name: subprocess_without_shell_equals_true
`

var expectedBanditOutput = [4]cqlinters.BanditMessage{
	{
		TestID: "B404", TestName: "blacklist", Confidence: "HIGH", Severity: "LOW",
		Filename: "./build/build/lib/mllint/cli.py", Line: 1,
		Text:        "Consider possible security implications associated with subprocess module.",
		MoreInfo:    "https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b404-import-subprocess",
		CodeSnippet: `1 import subprocess\n2 import sys\n3 import os\n`,
	},
	{
		TestID: "B603", TestName: "subprocess_without_shell_equals_true", Confidence: "HIGH", Severity: "LOW",
		Filename: "./build/build/lib/mllint/cli.py", Line: 9,
		Text:        "subprocess call - check for execution of untrusted input.",
		MoreInfo:    "https://bandit.readthedocs.io/en/latest/plugins/b603_subprocess_without_shell_equals_true.html",
		CodeSnippet: `8   os.chmod(mllint_exe, os.stat(mllint_exe).st_mode | 0o111) # Ensures mllint-exe is executable, equivalent to chmod +x\n9   return subprocess.run([mllint_exe] + sys.argv[1:]).returncode\n10 \n`,
	},
	{
		TestID: "B404", TestName: "blacklist", Confidence: "HIGH", Severity: "LOW",
		Filename: "./build/mllint/cli.py", Line: 1,
		Text:        "Consider possible security implications associated with subprocess module.",
		MoreInfo:    "https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b404-import-subprocess",
		CodeSnippet: `1 import subprocess\n2 import sys\n3 import os\n`,
	},
	{
		TestID: "B603", TestName: "subprocess_without_shell_equals_true", Confidence: "HIGH", Severity: "LOW",
		Filename: "./build/mllint/cli.py", Line: 10,
		Text:        "subprocess call - check for execution of untrusted input.",
		MoreInfo:    "https://bandit.readthedocs.io/en/latest/plugins/b603_subprocess_without_shell_equals_true.html",
		CodeSnippet: `9   os.chmod(mllint_exe, os.stat(mllint_exe).st_mode | 0o111) # Ensures mllint-exe is executable, equivalent to chmod +x\n10   return subprocess.run([mllint_exe] + sys.argv[1:]).returncode\n11 \n`,
	},
}
