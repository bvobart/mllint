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
