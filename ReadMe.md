# `mllint` — Linter for Machine Learning projects

`mllint` is a command-line utility to evaluate the quality of Machine Learning (ML) projects by means of static analysis of the project's repository. It measures the project's adherence to ML best practices, as collected and deduced from se4ml.github.io and Google's [Rules for ML](https://developers.google.com/machine-learning/guides/rules-of-ml).

TODO: write overview of linting rules

## Getting Started

`mllint` is published to PyPI, so it can be installed using `pip`:
```sh
pip install mllint
```

To run `mllint` in its default configuration, use one of the following commands:
```sh
# Run `mllint` on the project in the current folder
mllint

# Run `mllint` on the project in projects/my-ml-project
mllint projects/my-ml-project
```

To list all available or all enabled linting rules, use one of the following commands:
```sh
# List all available (implemented) linting rules
mllint list all

# List only the enabled rules for the project in the current folder.
mllint list enabled

# or for a project in projects/my-ml-project
mllint list enabled projects/my-ml-project
```

## Configuration

`mllint` can be configured using a `.mllint.yml` file that should be placed at the root of the project directory. This is a YAML file in which you can disable specific rules / linters, as well as configure specific settings for various linters.

An example `.mllint.yml` looks as follows:
```yaml
rules:
  disabled:
    - use-git # disables the 'use-git' linter
    - use-dependency-manager/single # disables the 'single' rule of the 'use-dependency-manager' linter.
    # - use-dependency-manager # this would disable the 'use-dependency-manager' linter and all of its rules entirely.
```

Here are some commands related to configuration:
```sh
# Print the configuration of the project in the current folder
mllint config

# Print the configuration of the project in projects/my-ml-project
mllint config projects/my-ml-project

# Print the default configuration (unless there's a folder called 'default' in the current dir)
mllint config default

# Create a valid `.mllint.yml` file with the default configuration
mllint config default -q > .mllint.yml
```

---

## Getting Started (development)

Clone this repository and run `go run .` in the root of this repository.
