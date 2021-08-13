# ML Project Report
**Project** | **Details**
--------|--------
Date    | Thu, 12 Aug 2021 16:00:23 +0200 
Path    | `/home/bart/tudelft/thesis/mllint-example-projects`
Config  | `pyproject.toml`
Default | No
Git: Remote URL | `git@github.com:bvobart/mllint-example-projects.git`
Git: Commit     | `3d559e7bb94d0a55714c67c77007a0b2eb124bb2`
Git: Branch     | `1-10-basic-tests`
Git: Dirty Workspace?  | Yes
Number of Python files | 7
Lines of Python code   | 197

---

## Reports

### Version Control (`version-control`) — **100.0**%

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project uses Git | `version-control/code/git`
✅ | 100.0% | 1 | Project should not have any large files in its Git history | `version-control/code/git-no-big-files`
✅ | 100.0% | 1 | DVC: Project uses Data Version Control | `version-control/data/dvc`
✅ | 100.0% | 1 | DVC: Is installed | `version-control/data/dvc-is-installed`
✅ | 100.0% | 1 | DVC: Folder '.dvc' should be committed to Git | `version-control/data/commit-dvc-folder`
✅ | 100.0% | 1 | DVC: Should have at least one remote data storage configured | `version-control/data/dvc-has-remote`
✅ | 100.0% | 1 | DVC: Should be tracking at least one data file | `version-control/data/dvc-has-files`
 | _Total_ | | | 
✅ | **100.0**% | | Version Control | `version-control`

### Dependency Management (`dependency-management`) — **100.0**%

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project properly keeps track of its dependencies | `dependency-management/use`
✅ | 100.0% | 1 | Project should only use one dependency manager | `dependency-management/single`
✅ | 100.0% | 1 | Project places its development dependencies in dev-dependencies | `dependency-management/use-dev`
 | _Total_ | | | 
✅ | **100.0**% | | Dependency Management | `dependency-management`

### Code Quality (`code-quality`) — **74.4**%

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project should use code quality linters | `code-quality/use-linters`
✅ | 100.0% | 1 | All code quality linters should be installed in the current environment | `code-quality/linters-installed`
❌ | 94.9% | 1 | Pylint reports no issues with this project | `code-quality/pylint/no-issues`
✅ | 100.0% | 1 | Pylint is configured for this project | `code-quality/pylint/is-configured`
✅ | 100.0% | 1 | Mypy reports no issues with this project | `code-quality/mypy/no-issues`
❌ | 0.0% | 1 | Black reports no issues with this project | `code-quality/black/no-issues`
✅ | 100.0% | 1 | isort reports no issues with this project | `code-quality/isort/no-issues`
✅ | 100.0% | 0 | isort is properly configured | `code-quality/isort/is-configured`
❌ | 0.0% | 1 | Bandit reports no issues with this project | `code-quality/bandit/no-issues`
 | _Total_ | | | 
❌ | **74.4**% | | Code Quality | `code-quality`

#### Details — Project should use code quality linters — ✅

Hooray, all linters detected:

- Pylint
- Mypy
- Black
- isort
- Bandit


#### Details — Pylint reports no issues with this project — ❌

Pylint reported **1** issues with your project:

- `tests/featurization_test.py:5,1` - _(W0511)_ TODO: implement tests for this module.


#### Details — Mypy reports no issues with this project — ✅

Congratulations, Mypy is happy with your project!

#### Details — Black reports no issues with this project — ❌

Black reported **2** files in your project that it would reformat:

- `tests/prepare_test.py`
- `src/mlbasic/prepare.py`

Black can fix these issues automatically when you run `black .` in your project.

#### Details — isort reports no issues with this project — ✅

Congratulations, `isort` is happy with your project!

#### Details — Bandit reports no issues with this project — ❌

Bandit reported **8** issues with your project:

- `tests/prepare_test.py:11` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:12` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:13` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:14` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:15` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:16` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:17` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)
- `tests/prepare_test.py:28` - _(B101, severity: LOW, confidence: HIGH)_ - Use of assert detected. The enclosed code will be removed when compiling to optimised byte code. [More Info](https://bandit.readthedocs.io/en/latest/plugins/b101_assert_used.html)


### Testing (`testing`) — **79.2**%

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project has automated tests | `testing/has-tests`
✅ | 100.0% | 1 | Project passes all of its automated tests | `testing/pass`
❌ | 17.0% | 1 | Project provides a test coverage report | `testing/coverage`
✅ | 100.0% | 1 | Tests should be placed in the tests folder | `testing/tests-folder`
 | _Total_ | | | 
❌ | **79.2**% | | Testing | `testing`

#### Details — Project has automated tests — ✅

Great! Your project contains **2** test files, which meets the minimum of **1** test files required.

This equates to **28.571429%** of Python files in your project being tests, which meets the target ratio of **20%**

#### Details — Project passes all of its automated tests — ✅

Congratulations, all **2** tests in your project passed!

#### Details — Project provides a test coverage report — ❌

Your project's tests achieved **13.6%** line test coverage, but **80.0%** is the target amount of test coverage to beat. You'll need to further improve your tests.

### Continuous Integration (`ci`) — **100.0**%

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project uses Continuous Integration (CI) | `ci/use`
 | _Total_ | | | 
✅ | **100.0**% | | Continuous Integration | `ci`

