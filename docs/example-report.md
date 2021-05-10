# ML Project Report
Project | Details
--------|--------
Path    | `/path/to/mllint-test-project`
Config  | `pyproject.toml`
Date    | Mon, 10 May 2021 23:19:48 +0200 
Number of Python files | 5
Lines of Python code | 204

---

## Reports

### Version Control (`version-control`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Code: Project uses Git | `version-control/code/git`
✅ | 100.0% | 1 | Code: Project should not use Git to track large files | `version-control/code/git-no-big-files`
✅ | 100.0% | 1 | Data: Project uses Data Version Control (DVC) | `version-control/data/dvc`
❌ | 0.0% | 1 | Data: DVC is installed | `version-control/data/dvc-is-installed`
✅ | 100.0% | 1 | Data: Folder '.dvc' should be committed to Git | `version-control/data/commit-dvc-folder`
❌ | 0.0% | 1 | Data: DVC should have at least one remote data storage configured | `version-control/data/dvc-has-remote`
❌ | 0.0% | 1 | Data: DVC should be tracking at least one data file | `version-control/data/dvc-has-files`

### Dependency Management (`dependency-management`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project properly keeps track of its dependencies | `dependency-management/use`
✅ | 100.0% | 1 | Project should only use one dependency manager | `dependency-management/single`

### Code Quality (`code-quality`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
❌ | 40.0% | 1 | Project should use code quality linters | `code-quality/use-linters`
✅ | 100.0% | 1 | All code quality linters should be installed in the current environment | `code-quality/linters-installed`
❌ | 0.0% | 1 | Pylint reports no issues with this project | `code-quality/pylint/no-issues`
❌ | 0.0% | 1 | Pylint is configured for this project | `code-quality/pylint/is-configured`
❌ | 11.8% | 1 | Mypy reports no issues with this project | `code-quality/mypy/no-issues`
❌ | 0.0% | 1 | Black reports no issues with this project | `code-quality/black/no-issues`
❌ | 0.0% | 1 | isort reports no issues with this project | `code-quality/isort/no-issues`
❌ | 0.0% | 0 | isort is properly configured | `code-quality/isort/is-configured`
❌ | 0.0% | 1 | Bandit reports no issues with this project | `code-quality/bandit/no-issues`

#### Details — Project should use code quality linters — ❌

Linters detected:

- Pylint
- Black


However, these linters were **missing** from your project:

- Mypy
- isort
- Bandit


We recommend that you start using these linters in your project to help you measure and maintain the quality of your code.

This rule will be satisfied, iff for each of these linters:
- **Either** there is a configuration file for this linter in the project
- **Or** the linter is a dependency of the project

Specifically, we recommend adding each linter to the development dependencies of your dependency manager,
e.g. using `poetry add --dev mypy` or `pipenv install --dev mypy`


#### Details — Pylint reports no issues with this project — ❌

Pylint reported **27** issues with your project:

- `src/evaluate.py:1,0` - _(C0114)_ Missing module docstring
- `src/evaluate.py:6,0` - _(E0401)_ Unable to import 'sklearn.metrics'
- `src/featurization.py:1,0` - _(C0114)_ Missing module docstring
- `src/featurization.py:3,0` - _(E0401)_ Unable to import 'pandas'
- `src/featurization.py:6,0` - _(E0401)_ Unable to import 'scipy.sparse'
- `src/featurization.py:9,0` - _(E0401)_ Unable to import 'sklearn.feature_extraction.text'
- `src/featurization.py:10,0` - _(E0401)_ Unable to import 'sklearn.feature_extraction.text'
- `src/featurization.py:32,0` - _(C0116)_ Missing function or method docstring
- `src/featurization.py:33,4` - _(C0103)_ Variable name "df" doesn't conform to snake_case naming style
- `src/featurization.py:44,0` - _(C0103)_ Argument name "df" doesn't conform to snake_case naming style
- `src/featurization.py:44,0` - _(C0116)_ Missing function or method docstring
- `src/featurization.py:53,31` - _(C0103)_ Variable name "fd" doesn't conform to snake_case naming style
- `src/featurization.py:55,4` - _(W0107)_ Unnecessary pass statement
- `src/featurization.py:5,0` - _(C0411)_ standard import "import pickle" should be placed before "import pandas as pd"
- `src/prepare.py:1,0` - _(C0114)_ Missing module docstring
- `src/prepare.py:20,0` - _(W0622)_ Redefining built-in 'input'
- `src/prepare.py:25,0` - _(C0116)_ Missing function or method docstring
- `src/prepare.py:25,18` - _(W0621)_ Redefining name 'fd_in' from outer scope (line 47)
- `src/prepare.py:25,25` - _(W0621)_ Redefining name 'fd_out_train' from outer scope (line 48)
- `src/prepare.py:25,39` - _(W0621)_ Redefining name 'fd_out_test' from outer scope (line 49)
- `src/prepare.py:41,15` - _(W0703)_ Catching too general exception Exception
- `src/train.py:1,0` - _(C0114)_ Missing module docstring
- `src/train.py:15,0` - _(W0622)_ Redefining built-in 'input'
- `src/train.py:6,0` - _(E0401)_ Unable to import 'sklearn.ensemble'
- `src/train_old.py:19,0` - _(W0622)_ Redefining built-in 'input'
- `src/train_old.py:10,0` - _(E0401)_ Unable to import 'sklearn.ensemble'
- `src/train_old.py:1,0` - _(R0801)_ Similar lines in 2 files
	```python
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
	```


#### Details — Mypy reports no issues with this project — ❌

Mypy reported **18** issues with your project:

- `src/evaluate.py:6` - Error: Cannot find implementation or library stub for module named 'sklearn.metrics'
- `src/evaluate.py:6` - Note: See https://mypy.readthedocs.io/en/latest/running_mypy.html#missing-imports
- `src/evaluate.py:6` - Error: Cannot find implementation or library stub for module named 'sklearn'
- `src/evaluate.py:37` - Error: Incompatible types in assignment (expression has type "TextIO", variable has type "BinaryIO")
- `src/evaluate.py:38` - Error: Argument 2 to "dump" has incompatible type "BinaryIO"; expected "IO[str]"
- `src/evaluate.py:40` - Error: Incompatible types in assignment (expression has type "TextIO", variable has type "BinaryIO")
- `src/evaluate.py:46` - Error: Argument 2 to "dump" has incompatible type "BinaryIO"; expected "IO[str]"
- `src/evaluate.py:48` - Error: Incompatible types in assignment (expression has type "TextIO", variable has type "BinaryIO")
- `src/evaluate.py:54` - Error: Argument 2 to "dump" has incompatible type "BinaryIO"; expected "IO[str]"
- `src/train_old.py:8` - Error: Skipping analyzing 'numpy': found module but no type hints or library stubs
- `src/train_old.py:10` - Error: Cannot find implementation or library stub for module named 'sklearn.ensemble'
- `src/train.py:4` - Error: Skipping analyzing 'numpy': found module but no type hints or library stubs
- `src/train.py:6` - Error: Cannot find implementation or library stub for module named 'sklearn.ensemble'
- `src/featurization.py:3` - Error: Cannot find implementation or library stub for module named 'pandas'
- `src/featurization.py:4` - Error: Skipping analyzing 'numpy': found module but no type hints or library stubs
- `src/featurization.py:6` - Error: Cannot find implementation or library stub for module named 'scipy.sparse'
- `src/featurization.py:6` - Error: Cannot find implementation or library stub for module named 'scipy'
- `src/featurization.py:9` - Error: Cannot find implementation or library stub for module named 'sklearn.feature_extraction.text'


#### Details — Black reports no issues with this project — ❌

Black reported **5** files in your project that it would reformat:

- `src/train_old.py`
- `src/train.py`
- `src/evaluate.py`
- `src/prepare.py`
- `src/featurization.py`

Black can fix these issues automatically when you run `black .` in your project.

#### Details — isort reports no issues with this project — ❌

isort reported **5** files in your project that it would fix:

- `src/evaluate.py` - Imports are incorrectly sorted and/or formatted.
- `src/train.py` - Imports are incorrectly sorted and/or formatted.
- `src/train_old.py` - Imports are incorrectly sorted and/or formatted.
- `src/prepare.py` - Imports are incorrectly sorted and/or formatted.
- `src/featurization.py` - Imports are incorrectly sorted and/or formatted.

isort can fix these issues automatically when you run `isort .` in your project.

#### Details — isort is properly configured — ❌

isort is not properly configured.
In order to be compatible with [Black](https://github.com/psf/black), which mllint also recommends using,
you should configure `isort` to use the `black` profile.
Furthermore, we recommend centralising your configuration in your `pyproject.toml`

Thus, ensure that your `pyproject.toml` contains at least the following section:

```toml
[tool.isort]
profile = "black"
```


#### Details — Bandit reports no issues with this project — ❌

Bandit reported **11** issues with your project:

- `src/evaluate.py:3` - _(B403, severity: LOW, confidence: HIGH)_ - Consider possible security implications associated with pickle module. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b403-import-pickle)
- `src/evaluate.py:20` - _(B301, severity: MEDIUM, confidence: HIGH)_ - Pickle and modules that wrap it can be unsafe when used to deserialize untrusted data, possible security issue. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b301-pickle)
- `src/evaluate.py:23` - _(B301, severity: MEDIUM, confidence: HIGH)_ - Pickle and modules that wrap it can be unsafe when used to deserialize untrusted data, possible security issue. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b301-pickle)
- `src/featurization.py:5` - _(B403, severity: LOW, confidence: HIGH)_ - Consider possible security implications associated with pickle module. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b403-import-pickle)
- `src/prepare.py:3` - _(B405, severity: LOW, confidence: HIGH)_ - Using xml.etree.ElementTree to parse untrusted XML data is known to be vulnerable to XML attacks. Replace xml.etree.ElementTree with the equivalent defusedxml package, or make sure defusedxml.defuse_stdlib() is called. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b405-import-xml-etree)
- `src/prepare.py:29` - _(B311, severity: LOW, confidence: HIGH)_ - Standard pseudo-random generators are not suitable for security/cryptographic purposes. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b311-random)
- `src/prepare.py:30` - _(B314, severity: MEDIUM, confidence: HIGH)_ - Using xml.etree.ElementTree.fromstring to parse untrusted XML data is known to be vulnerable to XML attacks. Replace xml.etree.ElementTree.fromstring with its defusedxml equivalent function or make sure defusedxml.defuse_stdlib() is called [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b313-b320-xml-bad-elementtree)
- `src/train.py:3` - _(B403, severity: LOW, confidence: HIGH)_ - Consider possible security implications associated with pickle module. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b403-import-pickle)
- `src/train.py:22` - _(B301, severity: MEDIUM, confidence: HIGH)_ - Pickle and modules that wrap it can be unsafe when used to deserialize untrusted data, possible security issue. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b301-pickle)
- `src/train_old.py:7` - _(B403, severity: LOW, confidence: HIGH)_ - Consider possible security implications associated with pickle module. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_imports.html#b403-import-pickle)
- `src/train_old.py:26` - _(B301, severity: MEDIUM, confidence: HIGH)_ - Pickle and modules that wrap it can be unsafe when used to deserialize untrusted data, possible security issue. [More Info](https://bandit.readthedocs.io/en/latest/blacklists/blacklist_calls.html#b301-pickle)


### Continuous Integration (`ci`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
❌ | 25.0% | 1 | Project uses Continuous Integration (CI) | `ci/use`

