# `mllint`'s model for evaluating ML project quality

TODO:
- [ ] Define aspect categories (e.g. version control, dependency management, data quality, code quality), based off SE4ML and Google's taxonomy.
- [ ] For each category, define the rules to be checked and how they can score the project.
- [ ] For each category, define how the scores of all rules amount to the score for this category.

## Score calculation

TODO:
- [ ] Define how the scores of all categories combine to create a score for the project.

Overall project score range: 0 to 100?

## Categories of evaluation

`mllint` evaluates ML projects from the perspectives outlined below.

### Version Control

Evaluates the use of version control software to keep track of changes in both code as well as data.

#### Rules

- use-git
- git-no-big-files
- use-dvc

### File and Folder structure

#### Rules

- check whether all data is placed inside a 'data' folder.
- I could take inspiration from  the file and folder structure produced by `mlapp`, MLFlow, `kedro`, and some of those other ML project generators / popular DS template projects.

### Dependency Management

Evaluates how the project's Python dependencies are managed

#### Rules

- all dependency management related rules

### Data Quality

Evaluates the quality of the data being fed into the ML application.

#### Rules

- check for incomplete data in dataset
- check whether data adheres to a specific schema? (i.e. are data types correct?)

### Code Quality

Evaluates the quality of the code that makes up the ML application.

#### Rules

- check whether the project uses linting and what linters are enabled.
- check which linting rules are enabled for each linter
- run Pylint with some custom configuration
- perhaps run mypy 

### Testing

Evaluates how well the project is being tested

#### Rules

- Measure test coverage
- Check which parts are being tested (ML part should also be tested)

### Continuous Integration (CI)

Evaluates whether the project uses Continuous Integration (CI) and what it is being used for.

#### Rules

- check whether the project has a CI config
- check whether linting / other CQ tooling is being used in CI?
- check whether tests are being run on CI

### Deployment

Evaluates to what degree the project is ready for deployment onto real-world infrastructure.

#### Rules

- has a Dockerfile? Dockerfile structure maybe?
- security compliance?

## Output

### Example console output

```
Linting project at  /home/bart/tudelft/thesis/mllint
No .mllint.yml or pyproject.toml found in project folder, using default configuration
---

Category:  Version Control (version-control)
Score:     100%

Passed Score Weight Rule                                                 Slug
✔️     100%  TODO   Code: Project uses Git                               code/git
✔️     100%         Code: Project isn't using Git to track large files   code/git-no-big-files
✔️     100%         Data: Project uses Data Version Control (DVC)        data/dvc
❌     0%           Data: DVC is installed                               data/dvc-is-installed
❌     0%           Data: .dvc folder should be committed to Git         data/dvc-commit-dvc-folder


---

Category:  File and Folder structure (file-and-folder-structure)
Score:     65%

Passed Score Weight Rule                                                 Slug
✔️     100%  TODO   Code: Project uses Git                               code/git
✔️     100%         Code: Project isn't using Git to track large files   code/git-no-big-files
✔️     100%         Data: Project uses Data Version Control (DVC)        data/dvc
❌     0%           Data: DVC is installed                               data/dvc-is-installed
❌     0%           Data: .dvc folder should be committed to Git         data/dvc-commit-dvc-folder


---
took: 13.37ms
```

### Example Markdown output

```md
# ML Project Report

Project | Details      
--------|------
Project | /home/bart/tudelft/thesis/mllint
Config  | default|configFilePath
Date    | $date
Commit  | $commitId

Overall Score             | 1337% | Weight
--------------------------|-------|--------
Version Control           | 80%   | TODO
File and Folder structure | 80%   | TODO

---

## Report per category

### Version Control (`version-control`) &mdash; 80%

Passed | Score | Weight | Rule                                               | Slug
:-----:|------:|--------|----------------------------------------------------|------
✔️     | 100%  | TODO   | Code: Project uses Git                             | code/git
✔️     | 100%  |        | Code: Project isn't using Git to track large files | code/git-no-big-files
✔️     | 100%  |        | Data: Project uses Data Version Control (DVC)      | data/dvc
❌     | 0%    |        | Data: DVC is installed                             | data/dvc-is-installed
❌     | 0%    |        | Data: .dvc folder should be committed to Git       | data/dvc-commit-dvc-folder

### File and Folder structure (`file-and-folder-structure`) &mdash; 80%

Passed | Score | Weight | Rule                                               | Slug
:-----:|------:|--------|----------------------------------------------------|------
✔️     | 100%  | TODO   | Code: Project uses Git                             | code/git
✔️     | 100%  |        | Code: Project isn't using Git to track large files | code/git-no-big-files
✔️     | 100%  |        | Data: Project uses Data Version Control (DVC)      | data/dvc
❌     | 0%    |        | Data: DVC is installed                             | data/dvc-is-installed
❌     | 0%    |        | Data: .dvc folder should be committed to Git       | data/dvc-commit-dvc-folder

---

## Appendices

Place any added appendices here, such as the output of Pylint or other tools.

```