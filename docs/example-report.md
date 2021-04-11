# ML Project Report
Project | Details
--------|--------
Path    | /path/to/mllint-test-project
Config  | pyproject.toml
Date    | Sun, 11 Apr 2021 20:16:25 +0200 

---

## Reports

### Version Control (`version-control`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Code: Project uses Git | version-control/code/git
✅ | 100.0% | 1 | Code: Project should not use Git to track large files | version-control/code/git-no-big-files
✅ | 100.0% | 1 | Data: Project uses Data Version Control (DVC) | version-control/data/dvc
❌ | 0.0% | 1 | Data: DVC is installed | version-control/data/dvc-is-installed
✅ | 100.0% | 1 | Data: Folder '.dvc' should be committed to Git | version-control/data/commit-dvc-folder
❌ | 0.0% | 1 | Data: DVC should have at least one remote data storage configured | version-control/data/dvc-has-remote
❌ | 0.0% | 1 | Data: DVC should be tracking at least one data file | version-control/data/dvc-has-files

### Dependency Management (`dependency-management`)

Passed | Score | Weight | Rule | Slug
:-----:|------:|-------:|------|-----
✅ | 100.0% | 1 | Project properly keeps track of its dependencies | dependency-management/use
✅ | 100.0% | 1 | Project should only use one dependency manager | dependency-management/single

