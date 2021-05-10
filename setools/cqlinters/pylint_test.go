package cqlinters_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/exec"
)

func TestPylint(t *testing.T) {
	l := cqlinters.Pylint{}
	require.Equal(t, cqlinters.TypePylint, l.Type())
	require.Equal(t, "Pylint", l.String())
	require.Equal(t, "pylint", l.DependencyName())

	exec.LookPath = func(file string) (string, error) { return "", errors.New("nope") }
	require.False(t, l.IsInstalled())
	exec.LookPath = func(file string) (string, error) { return "", nil }
	require.True(t, l.IsInstalled())
	exec.LookPath = exec.DefaultLookPath
}

const testPylintOutput = `
[
	{
		"type": "convention",
		"module": "train",
		"obj": "",
		"line": 1,
		"column": 0,
		"path": "src/train.py",
		"symbol": "missing-module-docstring",
		"message": "Missing module docstring",
		"message-id": "C0114"
	},
	{
		"type": "warning",
		"module": "train",
		"obj": "",
		"line": 15,
		"column": 0,
		"path": "src/train.py",
		"symbol": "redefined-builtin",
		"message": "Redefining built-in 'input'",
		"message-id": "W0622"
	},
	{
		"type": "warning",
		"module": "train_old",
		"obj": "",
		"line": 15,
		"column": 0,
		"path": "src/train_old.py",
		"symbol": "bad-indentation",
		"message": "Bad indentation. Found 4 spaces, expected 2",
		"message-id": "W0311"
	},
	{
		"type": "refactor",
		"module": "train_old",
		"obj": "",
		"line": 1,
		"column": 0,
		"path": "src/train_old.py",
		"symbol": "duplicate-code",
		"message": "Similar lines in 2 files\n==train:0\n==train_old:4\nimport sys\nimport os\nimport pickle\nimport numpy as np\nimport yaml\nfrom sklearn.ensemble import RandomForestClassifier\n\nparams = yaml.safe_load(open('params.yaml'))['train']\n\nif len(sys.argv) != 3:\n    sys.stderr.write('Arguments error. Usage:\\n')\n    sys.stderr.write('\\tpython train.py features model\\n')\n    sys.exit(1)\n\ninput = sys.argv[1]\noutput = sys.argv[2]\nseed = params['seed']\nn_est = params['n_est']\nmin_split = params['min_split']\n\nwith open(os.path.join(input, 'train.pkl'), 'rb') as fd:\n    matrix = pickle.load(fd)\n\nlabels = np.squeeze(matrix[:, 1].toarray())\nx = matrix[:, 2:]\n\nsys.stderr.write('Input matrix size {}\\n'.format(matrix.shape))\nsys.stderr.write('X matrix size {}\\n'.format(x.shape))\nsys.stderr.write('Y matrix size {}\\n'.format(labels.shape))\n\nclf = RandomForestClassifier(\n    n_estimators=n_est,\n    min_samples_split=min_split,\n    n_jobs=2,\n    random_state=seed\n)\n\nclf.fit(x, labels)\n\nwith open(output, 'wb') as fd:\n    pickle.dump(clf, fd)",
		"message-id": "R0801"
	}
]
`

var expectedPylintMessageStrings = [4]string{
	"`src/train.py[1,0]` - _(C0114)_ Missing module docstring",
	"`src/train.py[15,0]` - _(W0622)_ Redefining built-in 'input'",
	"`src/train_old.py[15,0]` - _(W0311)_ Bad indentation. Found 4 spaces, expected 2",
	fmt.Sprintf(
		`%s - _(R0801)_ Similar lines in 2 files
	%spython
	==train:0
	==train_old:4
	import sys
	import os
	import pickle
	import numpy as np
	import yaml
	from sklearn.ensemble import RandomForestClassifier
	
	params = yaml.safe_load(open('params.yaml'))['train']
	
	if len(sys.argv) != 3:
	    sys.stderr.write('Arguments error. Usage:\n')
	    sys.stderr.write('\tpython train.py features model\n')
	    sys.exit(1)
	
	input = sys.argv[1]
	output = sys.argv[2]
	seed = params['seed']
	n_est = params['n_est']
	min_split = params['min_split']
	
	with open(os.path.join(input, 'train.pkl'), 'rb') as fd:
	    matrix = pickle.load(fd)
	
	labels = np.squeeze(matrix[:, 1].toarray())
	x = matrix[:, 2:]
	
	sys.stderr.write('Input matrix size {}\n'.format(matrix.shape))
	sys.stderr.write('X matrix size {}\n'.format(x.shape))
	sys.stderr.write('Y matrix size {}\n'.format(labels.shape))
	
	clf = RandomForestClassifier(
	    n_estimators=n_est,
	    min_samples_split=min_split,
	    n_jobs=2,
	    random_state=seed
	)
	
	clf.fit(x, labels)
	
	with open(output, 'wb') as fd:
	    pickle.dump(clf, fd)
	%s`, "`src/train_old.py[1,0]`", "```", "```"),
}

func TestPylintRun(t *testing.T) {
	l := cqlinters.Pylint{}
	t.Run("EmptyProject", func(t *testing.T) {
		results, err := l.Run(api.Project{})
		require.NoError(t, err)
		require.Equal(t, []api.CQLinterResult{}, results)
	})

	t.Run("NormalProject+String", func(t *testing.T) {
		project := api.Project{
			Dir:         "test",
			PythonFiles: utils.Filenames{"file1", "file2", "file3"},
		}

		exec.CommandCombinedOutput = func(dir, name string, args ...string) ([]byte, error) {
			require.Equal(t, project.Dir, dir)
			require.Equal(t, "pylint", name)
			require.Equal(t, []string{"-f", "json", "file1", "file2", "file3"}, args)
			return []byte(testPylintOutput), errors.New("pylint always exits with an error when there are messages")
		}

		results, err := l.Run(project)
		require.NoError(t, err)
		require.Len(t, results, 4)
		for i, result := range results {
			require.IsType(t, cqlinters.PylintMessage{}, result)
			require.Equal(t, expectedPylintMessageStrings[i], result.String())
		}
	})
}
